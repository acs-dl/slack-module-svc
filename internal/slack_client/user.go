package slack_client

import (
	"fmt"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) UserFromApi(userId string) (*data.User, error) {
	user, err := s.superBotClient.GetUserInfo(userId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving users: %v.", user))
	}

	return &data.User{
		Username: &user.Name,
		Realname: &user.RealName,
		SlackId:  user.ID,
	}, nil
}
