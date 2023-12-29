package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetConversationsForUser(userId string) ([]slack.Channel, error) {
	channels, _, err := s.superBotClient.GetConversationsForUser(
		&slack.GetConversationsForUserParameters{
			UserID: userId,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversations for user", logan.F{
			"user_id": userId,
		})
	}

	return channels, nil
}
