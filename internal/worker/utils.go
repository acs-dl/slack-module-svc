package worker

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	slackGo "github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (w *worker) getUserFromDbBySlackId(slackId string) (*data.User, error) {
	usersQ := w.usersQ.New()
	user, err := usersQ.FilterBySlackIds(slackId).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db", logan.F{
			"slack_id": slackId,
		})
	}

	if user == nil {
		return nil, errors.From(errors.New("no such user in module"), logan.F{
			"slack_id": slackId,
		})
	}

	return user, nil
}

func (w *worker) getUsers() ([]slackGo.User, error) {
	return helpers.GetUsers(
		w.pqueues.BotPQueue,
		any(w.client.GetUsers),
		[]any{},
		pqueue.LowPriority,
	)
}

func (w *worker) getBillableInfo() (map[string]bool, error) {
	return helpers.GetBillableInfo(
		w.pqueues.UserPQueue,
		any(w.client.GetBillableInfo),
		pqueue.LowPriority,
	)
}

func (w *worker) getWorkspaceName() (string, error) {
	return helpers.GetWorkspaceName(
		w.pqueues.BotPQueue,
		any(w.client.GetWorkspaceName),
		[]any{},
		pqueue.LowPriority,
	)
}

func (w *worker) getConversationsForUser(userId string) ([]slackGo.Channel, error) {
	return helpers.GetConversationsForUser(
		w.pqueues.BotPQueue,
		any(w.client.GetConversationsForUser),
		[]interface{}{userId},
		pqueue.LowPriority,
	)
}

func (w *worker) getConversations() ([]data.Conversation, error) {
	return helpers.GetConversations(
		w.pqueues.BotPQueue,
		any(w.client.GetConversations),
		pqueue.LowPriority,
	)
}
