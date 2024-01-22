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

		conversations, nextCursor, err := c.getConversationsWrapper(&params, priority)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get conversations")
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

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return allConversations, nil
}

func (c *client) getConversationsWrapper(
	params *slack.GetConversationsParameters,
	priority int,
) ([]slack.Channel, string, error) {
	type response struct {
		conversations []slack.Channel
		nextCursor    string
	}

	var resp response
	err := doQueueRequest[response](QueueParameters{
		queue: c.pqueues.BotPQueue,
		function: func() (response, error) {
			conversations, nextCursor, err := c.botClient.GetConversations(params)
			return response{conversations, nextCursor}, err
		},
		args:     []any{},
		priority: priority,
	}, &resp)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get conversations", logan.F{
			"params": params,
		})
	}

	return resp.conversations, resp.nextCursor, nil
}
