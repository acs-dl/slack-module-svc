package slack_client

import (
	"context"
	"os"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
)

type ClientForSlack interface {
	UserFromApi(userId string) (*data.User, error)
	FetchUsers() ([]slack.User, error)
	WorkspaceName() (string, error)
	ConversationsForUser(userId string) ([]slack.Channel, error)
	BillableInfoForUser(userId string) (bool, error)

	ConversationFromApi(title string) ([]Conversation, error)
	//ConversationUserFromApi(user data.User, conversation Conversation) (*data.User, error)

	ConversationUsersFromApi(conversation Conversation) ([]data.User, error)
}

type slackStruct struct {
	log             *logan.Entry
	userAdminClient *slack.Client
	superBotClient  *slack.Client
}

type Conversation struct {
	Title         string
	Id            string
	MembersAmount int64
}

func NewSlack(cfg config.Config) ClientForSlack {
	return &slackStruct{
		log:             cfg.Log(),
		userAdminClient: slack.New(os.Getenv("USER_TOKEN")),
		superBotClient:  slack.New(os.Getenv("BOT_TOKEN")),
	}
}

func SlackClientInstance(ctx context.Context) ClientForSlack {
	return ctx.Value("slack").(ClientForSlack)
}
