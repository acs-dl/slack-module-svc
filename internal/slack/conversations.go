package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetConversations() ([]data.Conversation, error) {
	conversations, err := c.getConversations(func(_ slack.Channel) bool {
		return true
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversations")
	}

	return conversations, nil
}
