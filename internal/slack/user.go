package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUser(userId string, priority int) (*slack.User, error) {
	item, err := addFunctionInPQueue(c.pqueues.BotPQueue, c.botClient.GetUserInfo, []any{userId}, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	if err = item.Response.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user info from api", logan.F{
			"user_id": userId,
		})
	}

	user, ok := item.Response.Value.(*slack.User)
	if !ok {
		return nil, errors.New("failed to convert response")
	}

	return user, nil
}
