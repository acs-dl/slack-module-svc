package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUsers(priority int) ([]slack.User, error) {
	var users *[]slack.User
	err := doQueueRequest[*[]slack.User](QueueParameters{
		queue:    c.pqueues.BotPQueue,
		function: c.botClient.GetUsers,
		args:     []any{},
		priority: priority,
	}, &users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	return *users, nil
}
