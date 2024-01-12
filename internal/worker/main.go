package worker

import (
	"context"
	"time"

	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/processor"
	"github.com/acs-dl/slack-module-svc/internal/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
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
	runnerDelay   time.Duration
	estimatedTime time.Duration

	client  slack.Client
	pqueues *pqueue.PQueues
}

func New(cfg config.Config, ctx context.Context) Worker {
	return &worker{
		logger:        cfg.Log().WithField("runner", ServiceName),
		processor:     processor.ProcessorInstance(ctx),
		runnerDelay:   cfg.Runners().Worker,
		estimatedTime: time.Duration(0),

		client:  slack.New(cfg, ctx),
		pqueues: pqueue.PQueuesInstance(ctx),
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
	w.logger.Info("started processing permissions for all conversations")
	startTime := time.Now()
	conversations, err := w.client.GetConversations(pqueue.LowPriority)
	if err != nil {
		return errors.Wrap(err, "failed to get all conversations from slack api")
	}

	if err = w.processor.RefreshByConversations("from-worker", conversations...); err != nil {
		return errors.Wrap(err, "failed to refresh module")
	}

	w.estimatedTime = time.Since(startTime)
	w.logger.Info("finished processing permissions for all conversations")

	return nil
}

func (w *worker) RefreshModule() (string, error) {
	w.logger.Infof("started refreshing module")

	err := w.ProcessPermissions(context.Background())
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to refresh module")
	}

	w.logger.Infof("finished refreshing module")
	return data.SUCCESS, nil
}

func (w *worker) RefreshSubmodules(msg data.ModulePayload) (string, error) {
	w.logger.Infof("started refreshing submodules")

	if err := w.validateRefreshSubmodulesRequest(msg); err != nil {
		return data.FAILURE, errors.Wrap(err, "validation failed", logan.F{
			"links": msg.Links,
		})
	}

	var conversations []data.Conversation
	for _, link := range msg.Links {
		conversation, err := w.client.GetConversationsByLink(link, pqueue.LowPriority)
		if err != nil {
			return data.FAILURE, errors.Wrap(err, "failed to get conversation by link", logan.F{
				"link": link,
			})
		}

		conversations = append(conversations, conversation...)
	}

	if err := w.processor.RefreshByConversations("from-worker", conversations...); err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to refresh submodules")
	}

	w.logger.Infof("finished refreshsing submodules")
	return data.SUCCESS, nil
}

func (w *worker) GetEstimatedTime() time.Duration {
	return w.estimatedTime
}
