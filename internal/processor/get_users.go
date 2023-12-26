package processor

import (
	"fmt"
	"time"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/slack_client"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
		return errors.Wrap(err, "failed to get chat from api")
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
			if err := p.processUser(user, &msg, &workspaceName, &chat, &usersToUnverified); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to process user id:%s", user.SlackId))
			}
		}

		if err := p.sendUsers(msg.RequestId, usersToUnverified); err != nil {
			return errors.Wrap(err, "failed to publish users")
		}
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}

func (p *processor) getConversations(link string) ([]slack_client.Conversation, error) {
	return helpers.GetConversations(p.pqueues.SuperUserPQueue, any(p.client.ConversationFromApi), []any{any(link)}, pqueue.LowPriority)
}

func (p *processor) getUsersForChat(chat slack_client.Conversation) ([]data.User, error) {
	users, err := helpers.Users(p.pqueues.SuperUserPQueue, any(p.client.ConversationUsersFromApi), []any{any(chat)}, pqueue.LowPriority)
	if err != nil {
		return nil, err
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

func (p *processor) processUser(user data.User, msg *data.ModulePayload, workspaceName *string, chat *slack_client.Conversation, usersToUnverified *[]data.User) error {
	user.CreatedAt = time.Now()
	return p.managerQ.Transaction(func() error {
		if err := p.usersQ.Upsert(user); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to create user in db for message action with id `%s`", msg.RequestId))
		}

		dbUser, err := p.getUserFromDbBySlackId(user.SlackId)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get user from db for message action with id `%s`", msg.RequestId))
		}

		user.Id = dbUser.Id
		*usersToUnverified = append(*usersToUnverified, user)

		bill, err := helpers.GetBillableInfoForUser(p.pqueues.SuperUserPQueue, any(p.client.BillableInfoForUser), []interface{}{user.Id}, pqueue.LowPriority)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get billable info `%s`", msg.RequestId))
		}

		if err := p.permissionsQ.Upsert(data.Permission{
			RequestId:   msg.RequestId,
			WorkSpace:   *workspaceName,
			SlackId:     user.SlackId,
			Username:    *user.Username,
			AccessLevel: user.AccessLevel,
			Link:        msg.Link,
			CreatedAt:   user.CreatedAt,
			SubmoduleId: chat.Id,
			Bill:        bill,
		}); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to upsert permission in db for message action with id `%s`", msg.RequestId))
		}

		return nil
	})
}

func (p *processor) storeChatInDatabaseSafe(chat *slack_client.Conversation) error {
	err := p.conversationsQ.Upsert(data.Conversation{
		Title:         chat.Title,
		Id:            chat.Id,
		MembersAmount: chat.MembersAmount,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upsert chat")
	}

	return nil
}
