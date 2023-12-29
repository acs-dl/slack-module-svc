package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) ConversationsForUser(userId string) ([]slack.Channel, error) {

	channels, _, err := s.superBotClient.GetConversationsForUser(
		&slack.GetConversationsForUserParameters{
			UserID: userId,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get conversations for user %s.", userId))

	}

	return channels, nil
}
