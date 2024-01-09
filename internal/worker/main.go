package worker

import (
	"context"
	"time"

	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/processor"
	"github.com/acs-dl/slack-module-svc/internal/sender"
	"github.com/acs-dl/slack-module-svc/internal/slack"
	slackGo "github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
)

const (
	ServiceName       = data.ModuleName + "-worker"
	SetUsersAction    = "set_users"
	ProcessUserAction = "process_user"
)

type Worker interface {
	Run(ctx context.Context) error
	ProcessPermissions(_ context.Context) error
	GetEstimatedTime() time.Duration
	RefreshModule() (string, error)
	RefreshSubmodules(msg data.ModulePayload) (string, error)
}

type worker struct {
	logger        *logan.Entry
	processor     processor.Processor
	linksQ        data.Links
	permissionsQ  data.Permissions
	usersQ        data.Users
	runnerDelay   time.Duration
	estimatedTime time.Duration

	client          slack.Client
	pqueues         *pqueue.PQueues
	sender          sender.Sender
	unverifiedTopic string
}

func New(cfg config.Config, ctx context.Context) Worker {
	return &worker{
		logger:        cfg.Log().WithField("runner", ServiceName),
		processor:     processor.ProcessorInstance(ctx),
		linksQ:        postgres.NewLinksQ(cfg.DB()),
		permissionsQ:  postgres.NewPermissionsQ(cfg.DB()),
		usersQ:        postgres.NewUsersQ(cfg.DB()),
		runnerDelay:   cfg.Runners().Worker,
		estimatedTime: time.Duration(0),

		client:          slack.New(cfg),
		pqueues:         pqueue.PQueuesInstance(ctx),
		sender:          sender.SenderInstance(ctx),
		unverifiedTopic: cfg.Amqp().Unverified,
	}
}

func (w *worker) Run(ctx context.Context) error {
	running.WithBackOff(
		ctx,
		w.logger,
		ServiceName,
		w.ProcessPermissions,
		w.runnerDelay,
		w.runnerDelay,
		w.runnerDelay,
	)

	return nil
}

func (w *worker) ProcessPermissions(_ context.Context) error {
	w.logger.Info("started processing permissions")
	startTime := time.Now()

	err := w.processConversations()
	if err != nil {
		return errors.Wrap(err, "failed to process conversations")
	}

	usersToUnverified, err := w.processUnverifiedUsers()
	if err != nil {
		return errors.Wrap(err, "failed to process users")
	}

	err = w.sendUsers("", usersToUnverified)
	if err != nil {
		return errors.Wrap(err, "failed to publish unverified users")
	}

	w.estimatedTime = time.Since(startTime)
	w.logger.Info("finished processing permissions")

	return nil
}

func (w *worker) processUnverifiedUsers() ([]data.User, error) {
	w.logger.Info("getting users from api")
	users, err := w.getUsers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users from Slack API")
	}

	w.logger.Info("getting billable info from api")
	billableInfo, err := w.getBillableInfo()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get billable info from Slack API")
	}

	w.logger.Info("getting workspace name from api")
	workspaceName, err := w.getWorkspaceName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workspace name from Slack API")
	}

	usersToUnverified := make([]data.User, 0)
	for _, user := range users {
		userData, err := w.processUser(user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process user", logan.F{
				"slack_id": user.ID,
			})
		}

		usersToUnverified = append(usersToUnverified, *userData)

		w.logger.Info("inserting permissions into table 'permissions'")
		err = w.upsertPermissions(user, workspaceName, billableInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process permissions for user", logan.F{
				"user_id": user.ID,
			})
		}
	}

	return usersToUnverified, nil
}

func (w *worker) processUser(user slackGo.User) (*data.User, error) {
	w.logger.Info("inserting user into table 'users'")
	err := w.upsertUsers(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user into table 'users'")
	}

	dbUser, err := w.getUserFromDbBySlackId(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db", logan.F{
			"user_id": user.ID,
			"action":  ProcessUserAction,
		})
	}

	copiedName := user.Name
	copiedRealName := user.RealName

	return &data.User{
		Id:       dbUser.Id,
		Username: &copiedName,
		Realname: &copiedRealName,
		SlackId:  user.ID,
	}, nil
}

func (w *worker) processConversations() error {
	w.logger.Info("getting conversations from api")
	conversations, err := w.getConversations()
	if err != nil {
		return errors.Wrap(err, "failed to get conversations from Slack API")
	}

	if err := w.processor.UpsertConversations(conversations...); err != nil {
		return errors.Wrap(err, "failed to save conversation in db")
	}

	return nil
}

func (w *worker) upsertUsers(user slackGo.User) error {
	err := w.usersQ.Upsert(data.User{
		Username:  &user.Name,
		Realname:  &user.RealName,
		SlackId:   user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert user", logan.F{
			"user_id": user.ID,
		})
	}

	return nil
}

func (w *worker) upsertPermissions(user slackGo.User, workspaceName string, billableInfo map[string]bool) error {
	channels, err := w.getConversationsForUser(user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get user conversations")
	}

	for _, channel := range channels {
		bill, ok := billableInfo[user.ID]
		if !ok {
			return errors.From(errors.New("failed to get billable info for user"), logan.F{
				"user_id": user.ID,
			})
		}

		err := w.permissionsQ.Upsert(data.Permission{
			WorkSpace:   workspaceName,
			SlackId:     user.ID,
			Username:    user.Name,
			AccessLevel: w.userStatus(&user),
			Link:        channel.Name,
			SubmoduleId: channel.ID,
			Bill:        bill,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert a permission for user", logan.F{
				"user_id": user.ID,
			})
		}
	}
	return nil
}

func (w *worker) userStatus(user *slackGo.User) string {
	switch {
	case user.IsAdmin:
		return "admin"
	case user.IsOwner:
		return "owner"
	case user.IsPrimaryOwner:
		return "primary_owner"
	case user.IsStranger:
		return "stranger"
	case user.IsRestricted:
		return "restricted"
	default:
		return "member"
	}
}

func (w *worker) RefreshModule() (string, error) {
	w.logger.Infof("started refresh module")

	err := w.ProcessPermissions(context.Background())
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to refresh module")
	}

	w.logger.Infof("finished refresh module")
	return data.SUCCESS, nil
}

func (w *worker) RefreshSubmodules(msg data.ModulePayload) (string, error) {
	w.logger.Infof("started refresh submodules")

	for _, link := range msg.Links {
		w.logger.Infof("started refreshing `%s`", link)

		err := w.createPermissions(link)
		if err != nil {
			return data.FAILURE, errors.Wrap(err, "failed to create subs", logan.F{
				"link": link,
			})
		}
		w.logger.Infof("finished refreshing `%s`", link)
	}

	w.logger.Infof("finished refresh submodules")
	return data.SUCCESS, nil
}

func (w *worker) createPermissions(link string) error {
	if err := w.processor.HandleGetUsersAction(data.ModulePayload{
		RequestId: "from-worker",
		Link:      link,
	}); err != nil {
		return errors.Wrap(err, "failed to get users for link", logan.F{
			"link": link,
		})
	}

	return nil
}

func (w *worker) GetEstimatedTime() time.Duration {
	return w.estimatedTime
}
