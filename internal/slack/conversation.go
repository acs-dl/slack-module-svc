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

func (c *client) getConversationsWrapper(params slack.GetConversationsParameters, priority int) ([]slack.Channel, string, error) {
	type response struct {
		conversation []slack.Channel
		nextCursor   string
	}

	wrapperFunc := func() (response, error) {
		conversations, nextCursor, err := c.botClient.GetConversations(&params)
		return response{conversations, nextCursor}, err
	}

	item, err := addFunctionInPQueue(
		c.pqueues.BotPQueue,
		wrapperFunc,
		[]any{},
		priority,
	)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get conversations from api", logan.F{
			"params": params,
		})
	}

	result, ok := item.Response.Value.(response)
	if !ok {
		return nil, "", errors.New("failed to convert response")
	}

	return result.conversation, result.nextCursor, err
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

		conversations, nextCursor, err := c.getConversationsWrapper(params, priority)
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
