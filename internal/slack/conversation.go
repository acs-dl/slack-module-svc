package slack

import (
	"time"

	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetConversation(title string) ([]Conversation, error) {
	chats, err := s.findConversationByTitle(title)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find conversation by title", logan.F{
			"chat_title": title,
		})
	}

	return chats, nil
}

func (s *client) findConversationByTitle(title string) ([]Conversation, error) {

	//TODO: maybe use pq?
	var allConversations []Conversation
	limit := 20  // Maximum number of channels returned in one request
	cursor := "" // Used for pagination

	for {
		// Getting the list of channels with the current cursor
		params := slack.GetConversationsParameters{
			Limit:  limit,
			Cursor: cursor,
		}

		channels, nextCursor, err := s.botClient.GetConversations(&params)
		if err != nil {
			return nil, err
		}

		// Filter channels by name and add to the result
		for _, channel := range channels {
			if channel.Name == title {
				allConversations = append(allConversations, Conversation{
					Title:         channel.Name,
					Id:            channel.ID,
					MembersAmount: int64(channel.NumMembers),
				})
			}
		}

		// Check if there are more channels to be processed
		if nextCursor == "" {
			break
		}

		// Update the cursor for the next request
		cursor = nextCursor

		// Waiting to avoid exceeding the limit of requests per minute
		time.Sleep(3 * time.Second)
	}

	return allConversations, nil
}
