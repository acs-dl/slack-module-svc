package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUsers() ([]slack.User, error) {
	users, err := c.botClient.GetUsers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve users")
	}

	return users, nil
}
