package helpers

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/slack"
	slackGo "github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetConversations(queue *pqueue.PriorityQueue, function any, args []any, priority int) ([]slack.Conversation, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting chat from api")
	}

	conversations, ok := item.Response.Value.([]slack.Conversation)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return conversations, nil
}

func GetUsersWithChannels(queue *pqueue.PriorityQueue, function any, args []any, priority int) (map[string][]string, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting chat from api")
	}

	usersWithChannels, ok := item.Response.Value.(map[string][]string)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return usersWithChannels, nil
}

func GetUser(queue *pqueue.PriorityQueue, function any, args []any, priority int) (*data.User, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting user from api")
	}

	user, ok := item.Response.Value.(*data.User)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return user, nil
}

// func for worker
func GetUsers(queue *pqueue.PriorityQueue, function any, args []any, priority int) ([]slackGo.User, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting user from api")
	}

	user, ok := item.Response.Value.([]slackGo.User)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return user, nil
}

// func for worker
func Users(queue *pqueue.PriorityQueue, function any, args []any, priority int) ([]data.User, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting user from api")
	}

	user, ok := item.Response.Value.([]data.User)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return user, nil
}

func GetWorkspaceName(queue *pqueue.PriorityQueue, function any, args []any, priority int) (string, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return "", errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return "", errors.Wrap(err, "some error while getting chat from api")
	}

	workspaceName, ok := item.Response.Value.(string)
	if !ok {
		return "", errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return workspaceName, nil
}

func GetConversationsForUser(queue *pqueue.PriorityQueue, function any, args []any, priority int) ([]slackGo.Channel, error) {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting chat from api")
	}

	conversations, ok := item.Response.Value.([]slackGo.Channel)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting chat from api")
	}

	return conversations, nil
}

func GetBillableInfo(queue *pqueue.PriorityQueue, function any, priority int) (map[string]bool, error) {
	item, err := AddFunctionInPQueue(queue, function, []any{}, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "some error while getting billable info from api")
	}

	billableInfo, ok := item.Response.Value.(map[string]bool)
	if !ok {
		return nil, errors.Wrap(err, "wrong response type while getting billable info from api")
	}

	return billableInfo, nil
}

func RetrieveChat(chats []slack.Conversation, msg data.ModulePayload) *slack.Conversation {
	if len(chats) == 1 {
		return &chats[0]
	}

	for i := range chats {
		if chats[i].Title != msg.Link {
			continue
		}

		return &chats[i]
	}

	return nil
}

func GetRequestError(queue *pqueue.PriorityQueue, function any, args []any, priority int) error {
	item, err := AddFunctionInPQueue(queue, function, args, priority)
	if err != nil {
		return errors.Wrap(err, "failed to add function in pqueue")
	}

	return errors.Wrap(item.Response.Error, "response error occured")
}
