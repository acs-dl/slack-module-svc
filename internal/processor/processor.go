package processor

import (
	"context"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/manager"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/sender"
	"github.com/acs-dl/slack-module-svc/internal/slack_client"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	ServiceName = data.ModuleName + "-processor"

	UpdateSlackAction = "update_slack"

	SetUsersAction    = "set_users"
	DeleteUsersAction = "delete_users"

	RemoveSlackAction = "update_telegram"
)

type Processor interface {
	HandleGetUsersAction(msg data.ModulePayload) error
	HandleVerifyUserAction(msg data.ModulePayload) (string, error)
}

type processor struct {
	log             *logan.Entry
	client          slack_client.ClientForSlack
	permissionsQ    data.Permissions
	usersQ          data.Users
	conversationsQ  data.Conversations
	managerQ        *manager.Manager
	sender          *sender.Sender
	pqueues         *pqueue.PQueues
	unverifiedTopic string
	identityTopic   string
}

func NewProcessorAsInterface(cfg config.Config, ctx context.Context) interface{} {
	return interface{}(&processor{
		log:             cfg.Log().WithField("service", ServiceName),
		client:          slack_client.NewSlack(cfg),
		permissionsQ:    postgres.NewPermissionsQ(cfg.DB()),
		usersQ:          postgres.NewUsersQ(cfg.DB()),
		conversationsQ:  postgres.NewConversationsQ(cfg.DB()),
		managerQ:        manager.NewManager(cfg.DB()),
		sender:          sender.SenderInstance(ctx),
		pqueues:         pqueue.PQueuesInstance(ctx),
		unverifiedTopic: cfg.Amqp().Unverified,
		identityTopic:   cfg.Amqp().Identity,
	})
}

func ProcessorInstance(ctx context.Context) Processor {
	return ctx.Value(ServiceName).(Processor)
}

func CtxProcessorInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceName, entry)
}