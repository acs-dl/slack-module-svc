package processor

import (
	"time"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/slack"
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
		return p.handleValidationFailure(msg.RequestId, err)
	}

	chats, err := p.getConversations(msg.Link)
	if err != nil {
		return p.handleErrorWithMessageActionID(msg.RequestId, "failed to get chat from api", err)
	}

	for _, chat := range chats {
		if err := p.storeChatInDatabaseSafe(&chat); err != nil {
			return p.handleErrorWithMessageActionID(msg.RequestId, "failed to handle db chat flow", err)
		}

		users, err := p.getUsersForChat(chat)
		if err != nil {
			return p.handleErrorWithMessageActionID(msg.RequestId, "failed to get users from API", err)
		}

		if len(users) == 0 {
			p.log.Warnf("no user was found for message action with id `%s`", msg.RequestId)
			continue
		}

		workspaceName, err := p.getWorkspaceName()
		if err != nil {
			p.log.WithError(err).Error("failed to get workspaceName from API")
			return errors.Wrap(err, "failed to get workspaceName from API")
		}

		usersToUnverified := make([]data.User, 0)

		for _, user := range users {
			if err := p.processUser(user, &msg, &workspaceName, &chat, &usersToUnverified); err != nil {
				return err
			}
		}

		if err := p.sendUsers(msg.RequestId, usersToUnverified); err != nil {
			return p.handleErrorWithMessageActionID(msg.RequestId, "failed to publish users", err)
		}
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}

func (p *processor) handleValidationFailure(requestID string, err error) error {
	p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", requestID)
	return errors.Wrap(err, "failed to validate fields")
}

func (p *processor) handleErrorWithMessageActionID(requestID, message string, err error) error {
	p.log.WithError(err).Errorf("%s for message action with id `%s`", message, requestID)
	return errors.Wrap(err, message)
}

func (p *processor) getConversations(link string) ([]slack.Conversation, error) {
	return helpers.GetConversations(p.pqueues.SuperUserPQueue, any(p.client.ConversationFromApi), []any{any(link)}, pqueue.LowPriority)
}

func (p *processor) getUsersForChat(chat slack.Conversation) ([]data.User, error) {
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

func (p *processor) processUser(user data.User, msg *data.ModulePayload, workspaceName *string, chat *slack.Conversation, usersToUnverified *[]data.User) error {
	user.CreatedAt = time.Now()
	return p.managerQ.Transaction(func() error {
		if err := p.usersQ.Upsert(user); err != nil {
			p.log.WithError(err).Errorf("failed to create user in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to create user in user db")
		}

		dbUser, err := p.getUserFromDbBySlackId(user.SlackId)
		if err != nil {
			p.log.WithError(err).Errorf("failed to get user from db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to get user from")
		}

		user.Id = dbUser.Id
		*usersToUnverified = append(*usersToUnverified, user)

		// TODO: get billable info for the user
		//bill, err := helpers.GetBillableInfoForUser(p.pqueues.SuperUserPQueue, any(p.client.BillableInfoForUser), []interface{}{user.Id}, pqueue.LowPriority)

		if err := p.permissionsQ.Upsert(data.Permission{
			RequestId:   msg.RequestId,
			WorkSpace:   *workspaceName,
			SlackId:     user.SlackId,
			Username:    *user.Username,
			AccessLevel: user.AccessLevel,
			Link:        msg.Link,
			CreatedAt:   user.CreatedAt,
			SubmoduleId: chat.Id,
			Bill:        false,
		}); err != nil {
			p.log.WithError(err).Errorf("failed to upsert permission in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to upsert permission in db")
		}

		return nil
	})
}

func (p *processor) storeChatInDatabaseSafe(chat *slack.Conversation) error {
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
