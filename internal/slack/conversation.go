package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetConversationsByLink(title string, priority int) ([]data.Conversation, error) {
	conversations, err := c.getConversations(func(ch slack.Channel) bool {
		return ch.Name == title
	}, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find conversation by title", logan.F{
			"conversation_title": title,
		})
	}

	return conversations, nil
}

func (c *client) getConversations(predicate func(slack.Channel) bool, priority int) ([]data.Conversation, error) {
	var allConversations []data.Conversation
	limit := 20
	cursor := ""

	for {
		params := slack.GetConversationsParameters{
			Limit:  limit,
			Cursor: cursor,
		}

		response, err := c.paginationWrapper(func() (response, error) {
			conversations, nextCursor, err := c.botClient.GetConversations(&params)
			return response{conversations, nextCursor}, err
		}, priority)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get conversations")
		}

		conversations, ok := response.payload.([]slack.Channel)
		if !ok {
			return nil, errors.New("failed to convert response to conversations")
		}

		for _, conversation := range conversations {
			if predicate(conversation) {
				allConversations = append(allConversations, data.Conversation{
					Title:         conversation.Name,
					Id:            conversation.ID,
					MembersAmount: int64(conversation.NumMembers),
				})
			}
		}

		if response.nextCursor == "" {
			break
		}

		cursor = response.nextCursor
	}

	return allConversations, nil
}
