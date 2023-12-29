package slack

import (
	"context"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
)

type Client interface {
	UserFromApi(userId string) (*data.User, error)
	FetchUsers() ([]slack.User, error)
	WorkspaceName() (string, error)
	ConversationsForUser(userId string) ([]slack.Channel, error)
	BillableInfoForUser(userId string) (bool, error)
	GetBillableInfo() (map[string]slack.BillingActive, error)
	ConversationFromApi(title string) ([]Conversation, error)
	ConversationUsersFromApi(conversation Conversation) ([]data.User, error)
}

type client struct {
	log             *logan.Entry
	userAdminClient *slack.Client
	superBotClient  *slack.Client
}

type Conversation struct {
	Title         string
	Id            string
	MembersAmount int64
}

func New(cfg config.Config) Client {
	config := cfg.SlackParams()

	return &client{
		log:             cfg.Log(),
		userAdminClient: slack.New(config.UserToken),
		superBotClient:  slack.New(config.BotToken),
	}
}

func ClientInstance(ctx context.Context) Client {
	return ctx.Value("slack").(Client)
}
