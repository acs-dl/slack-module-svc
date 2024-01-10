package slack

import (
	"context"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
)

type Client interface {
	GetUser(userId string) (*data.User, error)
	GetUsers() ([]slack.User, error)
	GetWorkspaceName() (string, error)
	GetConversationsForUser(userId string) ([]slack.Channel, error)
	GetBillableInfoForUser(userId string) (bool, error)
	GetBillableInfo() (map[string]bool, error)
	GetConversationsByLink(title string) ([]data.Conversation, error)
	GetConversations() ([]data.Conversation, error)
	GetConversationUsers(conversation data.Conversation) ([]data.User, error)
}

type client struct {
	log        *logan.Entry
	userClient *slack.Client
	botClient  *slack.Client

	pqueues *pqueue.PQueues
}

func New(cfg config.Config, ctx context.Context) Client {
	config := cfg.SlackParams()

	return &client{
		log:        cfg.Log(),
		userClient: slack.New(config.UserToken),
		botClient:  slack.New(config.BotToken),
		pqueues:    pqueue.PQueuesInstance(ctx),
	}
}
