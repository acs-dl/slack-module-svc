package processor

import (
	"strconv"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateVerifyUser(msg data.ModulePayload) error {
	return validation.Errors{
		"user_id":  validation.Validate(msg.UserId, validation.Required),
		"username": validation.Validate(msg.Username, validation.Required),
	}.Filter()
}

func (p *processor) parseUserID(userID string) (int64, error) {
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		p.log.WithError(err).Errorf("failed to parse user id `%s`", userID)
		return 0, errors.Wrap(err, "failed to parse user id")
	}
	return userId, nil
}

func (p *processor) updateUserInDB(user *data.User, userID int64) error {
	user.Id = &userID
	if err := p.usersQ.Upsert(*user); err != nil {
		p.log.WithError(err).Errorf("failed to upsert user in db")
		return errors.Wrap(err, "failed to upsert user in db")
	}
	return nil
}

func (p *processor) HandleVerifyUserAction(msg data.ModulePayload) (string, error) {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	if err := p.validateVerifyUser(msg); err != nil {
		p.log.WithError(err).Errorf("failed to validate fields")
		return data.FAILURE, errors.Wrap(err, "failed to validate fields")
	}

	userId, err := p.parseUserID(msg.UserId)
	if err != nil {
		p.log.WithError(err).Errorf("failed to parse user id")
		return data.FAILURE, err
	}

	user, err := p.usersQ.FilterByUsername(msg.Username).Get()
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user by username from db")
		return data.FAILURE, err
	}

	if user == nil {
		p.log.Errorf("no user was found")
		return data.FAILURE, errors.New("no user was found")
	}

	user, err = p.getUserFromAPI(user.SlackId)
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user by id from Slack API")
		return data.FAILURE, err
	}

	if user == nil {
		p.log.Errorf("no user was found")
		return data.FAILURE, errors.New("no user was found")
	}

	if err := p.updateUserInDB(user, userId); err != nil {
		p.log.WithError(err).Errorf("failed to upsert user in db")
		return data.FAILURE, err
	}

	if err := p.sendUpdateUserSlack(msg.RequestId, data.ModulePayload{
		RequestId: msg.RequestId,
		UserId:    msg.UserId,
		Username:  msg.Username,
		Realname:  msg.Realname,
		Action:    UpdateSlackAction,
		SlackId:   msg.SlackId,
	}); err != nil {
		p.log.WithError(err).Errorf("failed to publish users")
		return data.FAILURE, errors.Wrap(err, "failed to publish users")
	}

	if err := p.SendDeleteUser(msg.RequestId, *user); err != nil {
		p.log.WithError(err).Errorf("failed to publish delete user")
		return data.FAILURE, errors.Wrap(err, "failed to publish delete user")
	}

	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return data.SUCCESS, nil
}

func (p *processor) getUserFromAPI(slackID string) (*data.User, error) {
	user, err := helpers.GetUser(p.pqueues.UserPQueue,
		any(p.client.UserFromApi),
		[]any{any(slackID)},
		pqueue.NormalPriority,
	)
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user from api")
		return nil, errors.Wrap(err, "failed to get user from api")
	}
	return user, nil
}
