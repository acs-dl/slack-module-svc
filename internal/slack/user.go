package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetUser(userId string) (*data.User, error) {
	user, err := c.botClient.GetUserInfo(userId)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving user", logan.F{
			"user_id": userId,
		})
	}

	return &data.User{
		Username: &user.Name,
		Realname: &user.RealName,
		SlackId:  user.ID,
	}, nil
}
