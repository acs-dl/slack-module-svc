package processor

import (
	"time"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	slackGo "github.com/acs-dl/slack-module-svc/internal/slack"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateGetUsers(msg data.ModulePayload) error {
	return validation.Errors{
		"link": validation.Validate(msg.Link, validation.Required),
	}.Filter()
}

func (p *processor) HandleGetUsersAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	if err := p.validateGetUsers(msg); err != nil {
		return errors.Wrap(err, "failed to validate user")
	}

	chats, err := p.getConversations(msg.Link)
	if err != nil {
		return errors.Wrap(err, "failed to get chat from api", logan.F{
			"link": msg.Link,
		})
	}

	billableInfo, err := p.getBillableInfo()
	if err != nil {
		return errors.Wrap(err, "failed to get billable info from Slack API")
	}

	for _, chat := range chats {
		if err := p.storeChatInDatabaseSafe(&chat); err != nil {
			return errors.Wrap(err, "failed to handle db chat flow")
		}

		users, err := p.getUsersForChat(chat)
		if err != nil {
			return errors.Wrap(err, "failed to get users from API")
		}

		if len(users) == 0 {
			p.log.Warnf("no user was found for message action with id `%s`", msg.RequestId)
			continue
		}

		workspaceName, err := p.getWorkspaceName()
		if err != nil {
			return errors.Wrap(err, "failed to get workspaceName from API")
		}

		usersToUnverified := make([]data.User, 0)
		for _, user := range users {
			if err := p.processUser(user, &msg, &workspaceName, &chat, &usersToUnverified, billableInfo); err != nil {
				return errors.Wrap(err, "failed to process user", logan.F{
					"slack_id": user.SlackId,
				})
			}
		}

		if err := p.sendUsers(msg.RequestId, usersToUnverified); err != nil {
			return errors.Wrap(err, "failed to publish users")
		}
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}

func (p *processor) getConversations(link string) ([]slackGo.Conversation, error) {
	return helpers.GetConversations(p.pqueues.SuperUserPQueue, any(p.client.ConversationFromApi), []any{any(link)}, pqueue.LowPriority)
}

func (p *processor) getBillableInfo() (map[string]slack.BillingActive, error) {
	return helpers.GetBillableInfo(p.pqueues.SuperUserPQueue, any(p.client.GetBillableInfo), pqueue.LowPriority)
}

func (p *processor) getUsersForChat(chat slackGo.Conversation) ([]data.User, error) {
	users, err := helpers.Users(p.pqueues.SuperUserPQueue, any(p.client.ConversationUsersFromApi), []any{any(chat)}, pqueue.LowPriority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users for chat", logan.F{
			"chat_id":    chat.Id,
			"chat_title": chat.Title,
		})
	}

	return users, nil
}

func (p *processor) getWorkspaceName() (string, error) {
	return helpers.WorkspaceName(
		p.pqueues.SuperUserPQueue,
		any(p.client.WorkspaceName),
		[]any{},
		pqueue.LowPriority,
	)
}

func (p *processor) processUser(
	user data.User,
	msg *data.ModulePayload,
	workspaceName *string,
	chat *slackGo.Conversation,
	usersToUnverified *[]data.User,
	billableInfo map[string]slack.BillingActive,
) error {
	user.CreatedAt = time.Now()
	return p.managerQ.Transaction(func() error {
		if err := p.managerQ.Users.Upsert(user); err != nil {
			return errors.Wrap(err, "failed to create user in db for message action", logan.F{
				"action": msg.RequestId,
			})
		}

		dbUser, err := p.getUserFromDbBySlackId(user.SlackId)
		if err != nil {
			return errors.Wrap(err, "failed to get user from db for message action", logan.F{
				"action": msg.RequestId,
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
			RequestId:   msg.RequestId,
			WorkSpace:   *workspaceName,
			SlackId:     user.SlackId,
			Username:    *user.Username,
			AccessLevel: user.AccessLevel,
			Link:        msg.Link,
			CreatedAt:   user.CreatedAt,
			SubmoduleId: chat.Id,
			Bill:        bill.BillingActive,
		}); err != nil {
			return errors.Wrap(err, "failed to upsert permission in db for message action", logan.F{
				"action": msg.RequestId,
			})
		}

		return nil
	})
}

func (p *processor) storeChatInDatabaseSafe(chat *slackGo.Conversation) error {
	err := p.managerQ.Conversations.Upsert(data.Conversation{
		Title:         chat.Title,
		Id:            chat.Id,
		MembersAmount: chat.MembersAmount,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upsert chat", logan.F{
			"chat_id":    chat.Id,
			"chat_title": chat.Title,
		})
	}

	return nil
}
