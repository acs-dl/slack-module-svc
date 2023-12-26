package slack_client

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) FetchUsers() ([]slack.User, error) {
	users, err := s.superBotClient.GetUsers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve users")
	}

	return users, nil
}
