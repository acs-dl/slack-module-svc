package slack

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetConversationUsers(conversation data.Conversation) ([]data.User, error) {
	users, err := c.getAllUsersFromConversation(conversation.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users from conversation")
	}

	return users, nil
}

func (c *client) getAllUsersFromConversation(conversationId string) ([]data.User, error) {

	//TODO: maybe youse priority queue?
	var users []data.User
	cursor := "" // For pagination

	for {
		// Get the list of users in the channel
		params := &slack.GetUsersInConversationParameters{
			ChannelID: conversationId,
			Cursor:    cursor,
		}
		userIDs, nextCursor, err := c.botClient.GetUsersInConversation(params)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get users in conversation")
		}

		// Getting information about each user
		for _, userID := range userIDs {
			user, err := c.botClient.GetUserInfo(userID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get user info")
			}
			users = append(users, data.User{
				Username:    &user.Name,
				Realname:    &user.RealName,
				SlackId:     user.ID,
				AccessLevel: c.userStatus(user),
			})
		}

		// Check if the next page of users is available
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
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
