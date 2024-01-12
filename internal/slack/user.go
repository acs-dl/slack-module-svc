package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUser(userId string, priority int) (*slack.User, error) {
	var user *slack.User
	err := doQueueRequest[*slack.User](QueueParameters{
		queue:    c.pqueues.BotPQueue,
		function: c.botClient.GetUserInfo,
		args:     []any{userId},
		priority: priority,
	}, &user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info", logan.F{
			"user_id": userId,
		})
	}

	return user, nil
}
