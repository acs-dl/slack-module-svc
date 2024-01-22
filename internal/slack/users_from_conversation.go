package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetConversationUsers(conversation data.Conversation, priority int) ([]data.User, error) {
	users, err := c.getAllUsersFromConversation(conversation.Id, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users from conversation")
	}

	return users, nil
}

func (c *client) getAllUsersFromConversation(conversationId string, priority int) ([]data.User, error) {
	var users []data.User
	cursor := ""

	for {
		params := slack.GetUsersInConversationParameters{
			ChannelID: conversationId,
			Cursor:    cursor,
		}

		userIDs, nextCursor, err := c.getUsersWrapper(&params, priority)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get users from conversation")
		}

		for _, userID := range userIDs {
			user, err := c.GetUser(userID, priority)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get user from api")
			}

			users = append(users, data.User{
				Username:    &user.Name,
				Realname:    &user.RealName,
				SlackId:     user.ID,
				AccessLevel: c.userStatus(user),
			})
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return users, nil
}

func (c *client) getUsersWrapper(
	params *slack.GetUsersInConversationParameters,
	priority int,
) ([]string, string, error) {
	type response struct {
		userIDs    []string
		nextCursor string
	}

	var resp response
	err := doQueueRequest[response](QueueParameters{
		queue: c.pqueues.BotPQueue,
		function: func() (response, error) {
			users, nextCursor, err := c.botClient.GetUsersInConversation(params)
			return response{users, nextCursor}, err
		},
		args:     []any{},
		priority: priority,
	}, &resp)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get users in conversation", logan.F{
			"params": params,
		})
	}

	return resp.userIDs, resp.nextCursor, nil
}

func (c *client) userStatus(user *slack.User) string {
	switch {
	case user.IsAdmin:
		return "admin"
	case user.IsOwner:
		return "owner"
	case user.IsPrimaryOwner:
		return "primary_owner"
	case user.IsStranger:
		return "stranger"
	case user.IsRestricted:
		return "restricted"
	default:
		return "member"
	}
}
