package processor

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) sendUsers(uuid string, users []data.User) error {

	unverifiedUsers := make([]data.UnverifiedUser, 0)
	for i := range users {
		if users[i].Id != nil {
			continue
		}
		permission, err := p.permissionsQ.FilterBySlackIds(users[i].SlackId).FilterByGreaterTime(users[i].CreatedAt).Get()

		if err != nil {
			p.log.WithError(err).Errorf("failed to select permissions by date `%s`", users[i].CreatedAt.String())
			return errors.Wrap(err, "failed to select permissions by date")
		}

		if permission == nil {
			continue
		}

		unverifiedUser := convertUserToUnverifiedUser(users[i], permission.Link)

		unverifiedUsers = append(unverifiedUsers, unverifiedUser)
	}

	marshaledPayload, err := json.Marshal(data.UnverifiedPayload{
		Action: SetUsersAction,
		Users:  unverifiedUsers,
	})
	if err != nil {
		p.log.WithError(err).Errorf("failed to marshal unverified users list")
		return errors.Wrap(err, "failed to marshal unverified users list")
	}

	err = p.sender.SendMessageToCustomChannel(p.unverifiedTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		p.log.WithError(err).Errorf("failed to publish users to `slack-module`")
		return errors.Wrap(err, "failed to publish users to `slack-module`")
	}

	p.log.Infof("successfully published users to `unverified-svc`")
	return nil
}

func (p *processor) SendDeleteUser(uuid string, user data.User) error {
	unverifiedUsers := make([]data.UnverifiedUser, 0)

	unverifiedUsers = append(unverifiedUsers, convertUserToUnverifiedUser(user, ""))

	marshaledPayload, err := json.Marshal(data.UnverifiedPayload{
		Action: DeleteUsersAction,
		Users:  unverifiedUsers,
	})
	if err != nil {
		p.log.WithError(err).Errorf("failed to marshal unverified users list")
		return errors.Wrap(err, "failed to marshal unverified users list")
	}

	err = p.sender.SendMessageToCustomChannel(p.unverifiedTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		p.log.WithError(err).Errorf("failed to publish users to `unverified-svc`")
		return errors.Wrap(err, "failed to publish users to `unverified-svc`")
	}

	p.log.Infof("successfully published users to `unverified-svc`")
	return nil
}

func (p *processor) buildMessage(uuid string, payload []byte) *message.Message {
	return &message.Message{
		UUID:     uuid,
		Metadata: nil,
		Payload:  payload,
	}
}

func convertUserToUnverifiedUser(user data.User, submodule string) data.UnverifiedUser {
	return data.UnverifiedUser{
		CreatedAt: user.CreatedAt,
		Module:    data.ModuleName,
		Submodule: submodule,
		ModuleId:  user.SlackId,
		Username:  user.Username,
		RealName:  user.Realname,
		SlackId:   user.SlackId,
	}
}

func (p *processor) sendUpdateUserSlack(uuid string, msg data.ModulePayload) error {
	marshaledPayload, err := json.Marshal(msg)
	if err != nil {
		p.log.WithError(err).Errorf("failed to marshal update slack info")
		return errors.Wrap(err, "failed to marshal update slack info")
	}

	err = p.sender.SendMessageToCustomChannel(p.identityTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		p.log.WithError(err).Errorf("failed to publish users to `identity-svc`")
		return errors.Wrap(err, "failed to publish users to `identity-svc`")
	}

	p.log.Infof("successfully published user to `identity-svc`")
	return nil
}
