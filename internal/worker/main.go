package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/processor"
	"github.com/acs-dl/slack-module-svc/internal/sender"
	"github.com/acs-dl/slack-module-svc/internal/slack_client"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
)

const (
	ServiceName    = data.ModuleName + "-worker"
	SetUsersAction = "set_users"
)

type IWorker interface {
	Run(ctx context.Context)
	ProcessPermissions(_ context.Context) error
	GetEstimatedTime() time.Duration
}

type Worker struct {
	logger        *logan.Entry
	processor     processor.Processor
	linksQ        data.Links
	permissionsQ  data.Permissions
	usersQ        data.Users
	runnerDelay   time.Duration
	estimatedTime time.Duration

	client          slack_client.ClientForSlack
	pqueues         *pqueue.PQueues
	sender          *sender.Sender
	unverifiedTopic string
}

func NewWorkerAsInterface(cfg config.Config, ctx context.Context) interface{} {
	return interface{}(&Worker{
		logger:        cfg.Log().WithField("runner", ServiceName),
		processor:     processor.ProcessorInstance(ctx),
		linksQ:        postgres.NewLinksQ(cfg.DB()),
		permissionsQ:  postgres.NewPermissionsQ(cfg.DB()),
		usersQ:        postgres.NewUsersQ(cfg.DB()),
		runnerDelay:   cfg.Runners().Worker,
		estimatedTime: time.Duration(0),

		client:          slack_client.NewSlack(cfg),
		pqueues:         pqueue.PQueuesInstance(ctx),
		sender:          sender.SenderInstance(ctx),
		unverifiedTopic: cfg.Amqp().Unverified,
	})
}

func (w *Worker) Run(ctx context.Context) error {
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

func (w *Worker) ProcessPermissions(_ context.Context) error {
	w.logger.Info("getting users from api")
	startTime := time.Now()

	usersStore, err := helpers.GetUsers(
		w.pqueues.SuperUserPQueue,
		any(w.client.FetchUsers),
		[]any{},
		pqueue.LowPriority,
	)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to get users from API ")
		return errors.Wrap(err, "some error while getting users from api")
	}

	w.logger.Info("getting workspaceName from api")
	workspaceName, err := helpers.WorkspaceName(
		w.pqueues.SuperUserPQueue,
		any(w.client.WorkspaceName),
		[]any{},
		pqueue.LowPriority,
	)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to get workspaceName from API ")
		return errors.Wrap(err, "some error while getting workspaceName from api")
	}

	usersToUnverified := make([]data.User, 0)
	for _, user := range usersStore {

		w.logger.Info("inserting user into table 'users' from db")
		err := w.retrieveAndUpsertUsers(user)
		if err != nil {
			w.logger.WithError(err).Errorf("failed to insert users into table 'users' from db ")
			return errors.Wrap(err, "some error while inserting users into table 'users' from db")
		}

		msg := data.ModulePayload{
			RequestId: "",
			Link:      "",
		}

		dbUser, err := w.getUserFromDbBySlackId(user.ID)
		if err != nil {
			w.logger.WithError(err).Errorf("failed to get user from db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to get user from")
		}

		copiedName := user.Name
		copiedRealName := user.RealName

		userData := data.User{
			Id:       dbUser.Id,
			Username: &copiedName,
			Realname: &copiedRealName,
			SlackId:  user.ID,
		}

		usersToUnverified = append(usersToUnverified, userData)
		fmt.Printf("len %d\n", len(usersToUnverified))

		w.logger.Info("inserting permissions into table 'permissions' from db")
		err = w.retrieveAndUpsertPermissions(user, workspaceName)
		if err != nil {
			w.logger.WithError(err).Errorf("failed to process permissions for user")
			return errors.Wrap(err, "some error while processing user permissions")
		}

	}

	msg := data.ModulePayload{}

	err = w.sendUsers(msg.RequestId, usersToUnverified)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to publish users for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to publish users")
	}

	w.logger.Info("ProcessPermissions completed")

	w.estimatedTime = time.Since(startTime)
	return nil
}

func (w *Worker) retrieveAndUpsertUsers(user slack.User) error {
	err := w.usersQ.Upsert(data.User{
		Username:  &user.Name,
		Realname:  &user.RealName,
		SlackId:   user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		logan.New().WithError(err).Error("error insert user")
		return err

	}

	return nil
}

func (w *Worker) retrieveAndUpsertPermissions(user slack.User, workspaceName string) error {
	channels, err := helpers.GetConversationsForUser(
		w.pqueues.SuperUserPQueue,
		any(w.client.ConversationsForUser),
		[]interface{}{user.ID},
		pqueue.LowPriority,
	)
	if err != nil {
		return errors.Wrap(err, "Error getting user conversations.")
	}

	//TODO: uncomment when billable info will be useful for you
	//bill, err := helpers.GetBillableInfoForUser(w.pqueues.SuperUserPQueue, any(w.client.BillableInfoForUser), []interface{}{user.ID}, pqueue.LowPriority)

	for _, channel := range channels {

		err := w.permissionsQ.Upsert(data.Permission{
			WorkSpace:   workspaceName,
			SlackId:     user.ID,
			Username:    user.Name,
			AccessLevel: w.userStatus(&user),
			Link:        channel.Name,
			SubmoduleId: channel.ID,
			Bill:        false, //TODO: just use 'bill' here
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		if err != nil {
			return errors.Wrap(err, "Error insert permission")
		}
	}
	return nil
}

func (w *Worker) userStatus(user *slack.User) string {
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

func (w *Worker) RefreshModule() (string, error) {
	w.logger.Infof("started refresh module")

	err := w.ProcessPermissions(context.Background())
	if err != nil {
		return data.FAILURE, err
	}

	w.logger.Infof("finished refresh module")
	return data.SUCCESS, nil
}

func (w *Worker) RefreshSubmodules(msg data.ModulePayload) (string, error) {
	w.logger.Infof("started refresh submodules")

	for _, link := range msg.Links {
		w.logger.Infof("started refreshing `%s`", link)

		err := w.createPermissions(link)
		if err != nil {
			w.logger.Infof("failed to create subs for link `%s", link)
			return data.FAILURE, errors.Wrap(err, "failed to create subs")
		}
		w.logger.Infof("finished refreshing `%s`", link)
	}

	w.logger.Infof("finished refresh submodules")
	return data.SUCCESS, nil
}

func (w *Worker) createPermissions(link string) error {
	if err := w.processor.HandleGetUsersAction(data.ModulePayload{
		RequestId: "from-worker",
		Link:      link,
	}); err != nil {
		w.logger.Infof("failed to get users for link `%s`", link)
		return errors.Wrap(err, "failed to get users")
	}

	return nil
}

func (w *Worker) GetEstimatedTime() time.Duration {
	return w.estimatedTime
}
