package processor

import (
	"time"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) RefreshByConversations(requestId string, conversations ...data.Conversation) error {
	billableInfo, err := p.client.GetBillableInfo(pqueue.LowPriority)
	if err != nil {
		return errors.Wrap(err, "failed to get billable info from Slack API")
	}

	workspaceName, err := p.client.GetWorkspaceName(pqueue.LowPriority)
	if err != nil {
		return errors.Wrap(err, "failed to get workspaceName from API")
	}

	for _, conversation := range conversations {
		err := p.processConversation(conversation, requestId, workspaceName, billableInfo)
		if err != nil {
			return errors.Wrap(err, "failed to process conversation")
		}
	}

	p.log.Infof("finished handling message with %s action", requestId)
	return nil
}

func (p *processor) processConversation(
	conversation data.Conversation,
	requestId,
	workspaceName string,
	billableInfo map[string]bool,
) error {
	if err := p.managerQ.Conversations.Upsert(conversation); err != nil {
		return errors.Wrap(err, "failed to upsert conversation", logan.F{
			"id":    conversation.Id,
			"title": conversation.Title,
		})
	}

	if conversation.MembersAmount == 0 {
		p.log.Warnf("no users were found in conversation %s", conversation.Title)
		return nil
	}

	users, err := p.client.GetConversationUsers(conversation, pqueue.LowPriority)
	if err != nil {
		return errors.Wrap(err, "failed to get users from conversation", logan.F{
			"id":    conversation.Id,
			"title": conversation.Title,
		})
	}

	usersToUnverified := make([]data.User, 0)
	for _, user := range users {
		bill, ok := billableInfo[user.SlackId]
		if !ok {
			return errors.From(errors.New("failed to get billable info for user"), logan.F{
				"slack_id": user.SlackId,
			})
		}

		user.CreatedAt = time.Now()
		permission := data.Permission{
			RequestId:   requestId,
			WorkSpace:   workspaceName,
			SlackId:     user.SlackId,
			Username:    *user.Username,
			AccessLevel: user.AccessLevel,
			Link:        conversation.Title,
			CreatedAt:   user.CreatedAt,
			SubmoduleId: conversation.Id,
			Bill:        bill,
		}

		if err := p.processUser(&user, permission); err != nil {
			return errors.Wrap(err, "failed to process user", logan.F{
				"slack_id": user.SlackId,
				"action":   requestId,
			})
		}
		usersToUnverified = append(usersToUnverified, user)
	}

	if err := p.sendUsers(requestId, usersToUnverified); err != nil {
		return errors.Wrap(err, "failed to publish users")
	}

	return nil
}

func (p *processor) processUser(
	user *data.User,
	permission data.Permission,
) error {
	err := p.managerQ.Transaction(func() error {
		resultId, err := p.managerQ.Users.Upsert(*user)
		if err != nil {
			return errors.Wrap(err, "failed to upsert user")
		}

		if err := p.managerQ.Permissions.Upsert(permission); err != nil {
			return errors.Wrap(err, "failed to upsert permission")
		}

		user.Id = resultId
		return nil
	})

	return errors.Wrap(err, "transaction failed")
}
