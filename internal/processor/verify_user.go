package processor

import (
	"strconv"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3"
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
		return 0, errors.Wrap(err, "failed to parse user", logan.F{
			"user_id": userID,
		})
	}
	return userId, nil
}

func (p *processor) updateUserInDB(user *data.User, userID int64) error {
	user.Id = &userID
	if err := p.managerQ.Users.Upsert(*user); err != nil {
		return errors.Wrap(err, "failed to upsert user in db", logan.F{
			"user_id": userID,
		})
	}
	return nil
}

func (p *processor) HandleVerifyUserAction(msg data.ModulePayload) (string, error) {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	if err := p.validateVerifyUser(msg); err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to validate fields")
	}

	userId, err := p.parseUserID(msg.UserId)
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to parse user id", logan.F{
			"user_id": msg.UserId,
		})
	}

	user, err := p.managerQ.Users.FilterByUsername(msg.Username).Get()
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to get user by username", logan.F{
			"username": msg.Username,
		})
	}

	if user == nil {
		return data.FAILURE, errors.New("no user was found")
	}

	user, err = p.getUser(user.SlackId)
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to get user by id from Slack API", logan.F{
			"slack_id": user.SlackId,
		})
	}

	if user == nil {
		return data.FAILURE, errors.New("no user was found")
	}

	if err := p.updateUserInDB(user, userId); err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to upsert user in db", logan.F{
			"user_id": userId,
		})
	}

	if err := p.sendUpdateUserSlack(msg.RequestId, data.ModulePayload{
		RequestId: msg.RequestId,
		UserId:    msg.UserId,
		Username:  msg.Username,
		Realname:  msg.Realname,
		Action:    UpdateSlackAction,
		SlackId:   msg.SlackId,
	}); err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to publish users")
	}

	if err := p.SendDeleteUser(msg.RequestId, *user); err != nil {
		return data.FAILURE, errors.Wrap(err, "failed to publish delete user")
	}

	p.log.Infof("finish handle message action `%s`", msg.RequestId)
	return data.SUCCESS, nil
}

func (p *processor) getUser(slackID string) (*data.User, error) {
	user, err := helpers.GetUser(p.pqueues.BotPQueue,
		any(p.client.GetUser),
		[]any{any(slackID)},
		pqueue.NormalPriority,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from api", logan.F{
			"slack_id": slackID,
		})
	}
	return user, nil
}
