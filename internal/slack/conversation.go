package slack

import (
	"time"

	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetConversationsByLink(title string) ([]Conversation, error) {
	chats, err := s.getConversations(func(ch slack.Channel) bool {
		return ch.Name == title
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find conversation by title", logan.F{
			"chat_title": title,
		})
	}

	return chats, nil
}

func (s *client) getConversations(predicate func(slack.Channel) bool) ([]Conversation, error) {
	// TODO: consider creating a wrapper to use a pqueue
	var allConversations []Conversation
	limit := 20
	cursor := ""

	for {
		params := slack.GetConversationsParameters{
			Limit:  limit,
			Cursor: cursor,
		}

		channels, nextCursor, err := s.botClient.GetConversations(&params)
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			if predicate(channel) {
				allConversations = append(allConversations, Conversation{
					Title:         channel.Name,
					Id:            channel.ID,
					MembersAmount: int64(channel.NumMembers),
				})
			}
		}

		if nextCursor == "" {
			break
		}

		cursor = nextCursor

		// Waiting to avoid exceeding the limit of requests per minute
		time.Sleep(3 * time.Second)
	}

	return allConversations, nil
}
