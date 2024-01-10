package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUsers(priority int) ([]slack.User, error) {
	item, err := addFunctionInPQueue(c.pqueues.BotPQueue, c.botClient.GetUsers, []any{}, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	if err = item.Response.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get users from api")
	}

	users, ok := item.Response.Value.([]slack.User)
	if !ok {
		return nil, errors.New("failed to convert response")
	}

	return users, nil
}
