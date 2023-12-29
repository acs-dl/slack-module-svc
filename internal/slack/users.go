package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) FetchUsers() ([]slack.User, error) {

	users, err := s.superBotClient.GetUsers()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving users: %v.", users))
	}

	return users, nil
}
