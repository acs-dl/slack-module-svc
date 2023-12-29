package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
)

type Client interface {
	GetUser(userId string) (*data.User, error)
	GetUsers() ([]slack.User, error)
	GetWorkspaceName() (string, error)
	GetConversationsForUser(userId string) ([]slack.Channel, error)
	GetBillableInfoForUser(userId string) (bool, error)
	GetBillableInfo() (map[string]slack.BillingActive, error)
	GetConversation(title string) ([]Conversation, error)
	GetConversationUsers(conversation Conversation) ([]data.User, error)
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
