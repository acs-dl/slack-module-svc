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
		if err := p.processUser(user, requestId, workspaceName, conversation, &usersToUnverified, billableInfo); err != nil {
			return errors.Wrap(err, "failed to process user", logan.F{
				"slack_id": user.SlackId,
			})
		}
	}

	if err := p.sendUsers(requestId, usersToUnverified); err != nil {
		return errors.Wrap(err, "failed to publish users")
	}

	return nil
}

func (p *processor) processUser(
	user data.User,
	requestId string,
	workspaceName string,
	conversation data.Conversation,
	usersToUnverified *[]data.User,
	billableInfo map[string]bool,
) error {
	user.CreatedAt = time.Now()
	return p.managerQ.Transaction(func() error {
		if err := p.managerQ.Users.Upsert(user); err != nil {
			return errors.Wrap(err, "failed to upsert user", logan.F{
				"action":   requestId,
				"slack_id": user.SlackId,
			})
		}

		dbUser, err := p.getUserFromDbBySlackId(user.SlackId)
		if err != nil {
			return errors.Wrap(err, "failed to get user from db for message action", logan.F{
				"action":   requestId,
				"slack_id": user.SlackId,
			})
		}

		user.Id = dbUser.Id
		*usersToUnverified = append(*usersToUnverified, user)

		bill, ok := billableInfo[user.SlackId]
		if !ok {
			return errors.From(errors.New("failed to get billable info for user"), logan.F{
				"slack_id": user.SlackId,
			})
		}

		if err := p.managerQ.Permissions.Upsert(data.Permission{
			RequestId:   requestId,
			WorkSpace:   workspaceName,
			SlackId:     user.SlackId,
			Username:    *user.Username,
			AccessLevel: user.AccessLevel,
			Link:        conversation.Title,
			CreatedAt:   user.CreatedAt,
			SubmoduleId: conversation.Id,
			Bill:        bill,
		}); err != nil {
			return errors.Wrap(err, "failed to upsert permission", logan.F{
				"action":   requestId,
				"slack_id": user.SlackId,
			})
		}

		return nil
	})
}
