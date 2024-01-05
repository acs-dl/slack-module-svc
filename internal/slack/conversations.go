package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetConversations() ([]Conversation, error) {
	var allConversations []Conversation
	limit := 100
	cursor := ""

	for {
		channels, nextCursor, err := s.botClient.GetConversations(&slack.GetConversationsParameters{
			Limit:  limit,
			Cursor: cursor,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get conversations from slack api")
		}

		for _, channel := range channels {
			allConversations = append(allConversations, Conversation{
				Title:         channel.Name,
				Id:            channel.ID,
				MembersAmount: int64(channel.NumMembers),
			})
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return allConversations, nil
}
