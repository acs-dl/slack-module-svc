package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetConversations() ([]Conversation, error) {
	chats, err := s.getConversations(func(_ slack.Channel) bool {
		return true
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversations")
	}

	return chats, nil
}
