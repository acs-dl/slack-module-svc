package slack_client

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) ConversationFromApi(title string) ([]Conversation, error) {
	chats, err := s.getConversationFlow(title)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get chat `%s`", title))
	}

	return chats, nil
}

func (s *slackStruct) getConversationFlow(title string) ([]Conversation, error) {
	chats, err := s.findConversationByTitle(title)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *slackStruct) findConversationByTitle(title string) ([]Conversation, error) {

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

		channels, nextCursor, err := s.superBotClient.GetConversations(&params)
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
