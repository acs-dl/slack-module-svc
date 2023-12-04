package slack_client

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) ConversationUsersFromApi(conversation Conversation) ([]data.User, error) {

	users, err := s.getChatMembers(conversation)
	if err != nil {
		s.log.Errorf("failed to get chat members")
		return nil, errors.Wrap(err, "failed to get chat members")
	}

	return users, nil
}

func (s *slackStruct) getChatMembers(conversation Conversation) ([]data.User, error) {
	users, err := s.getAllUsers(conversation.Id)
	if err != nil {
		s.log.Errorf("failed to get all users")
		return nil, err
	}

	return users, nil
}

func (s *slackStruct) getAllUsers(id string) ([]data.User, error) {
	users := make([]data.User, 0)
	var err error = nil

	users, err = s.getAllUsersFromConversation(id)
	if err != nil {
		s.log.Errorf("failed to get all users from chat")
		return nil, err
	}

	s.log.Infof("found `%d` users", len(users))
	return users, nil
}

func (s *slackStruct) getAllUsersFromConversation(chatId string) ([]data.User, error) {

	//TODO: maybe youse priority queue?
	var users []data.User
	cursor := "" // For pagination

	for {
		// Get the list of users in the channel
		params := &slack.GetUsersInConversationParameters{
			ChannelID: chatId,
			Cursor:    cursor,
		}
		userIDs, nextCursor, err := s.superBotClient.GetUsersInConversation(params)
		if err != nil {
			return nil, err
		}

		// Getting information about each user
		for _, userID := range userIDs {
			user, err := s.superBotClient.GetUserInfo(userID)
			if err != nil {
				return nil, err
			}
			users = append(users, data.User{
				Username:    &user.Name,
				Realname:    &user.RealName,
				SlackId:     user.ID,
				AccessLevel: s.userStatus(user),
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

func (s *slackStruct) userStatus(user *slack.User) string {
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
