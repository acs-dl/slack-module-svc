package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
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

		response, err := c.paginationWrapper(func() (response, error) {
			users, nextCursor, err := c.botClient.GetUsersInConversation(&params)
			return response{users, nextCursor}, err
		}, priority)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get users in conversation")
		}

		userIDs, ok := response.payload.([]string)
		if !ok {
			return nil, errors.New("failed to convert response to slice of users")
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

		if response.nextCursor == "" {
			break
		}
		cursor = response.nextCursor
	}

	return users, nil
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
