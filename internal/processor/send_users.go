package processor

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) sendUsers(uuid string, users []data.User) error {
	unverifiedUsers := make([]data.UnverifiedUser, 0)
	for i := range users {
		if users[i].Id != nil {
			continue
		}

		permission, err := p.managerQ.Permissions.FilterBySlackIds(users[i].SlackId).FilterByGreaterTime(users[i].CreatedAt).Get()
		if err != nil {
			return errors.Wrap(err, "failed to select permissions by date", logan.F{
				"date": users[i].CreatedAt.String(),
			})
		}

		if permission == nil {
			continue
		}

		unverifiedUsers = append(unverifiedUsers, data.ConvertUserToUnverifiedUser(users[i], permission.Link))
	}

	marshaledPayload, err := json.Marshal(data.UnverifiedPayload{
		Action: SetUsersAction,
		Users:  unverifiedUsers,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal unverified users list")
	}

	err = p.sender.SendMessageToCustomChannel(p.unverifiedTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		return errors.Wrap(err, "failed to publish users", logan.F{
			"topic": p.unverifiedTopic,
		})
	}

	p.log.Infof("successfully published users to `%s`", p.unverifiedTopic)
	return nil
}

func (p *processor) SendDeleteUser(uuid string, user data.User) error {
	unverifiedUsers := make([]data.UnverifiedUser, 0)
	unverifiedUsers = append(unverifiedUsers, data.ConvertUserToUnverifiedUser(user, ""))
	marshaledPayload, err := json.Marshal(data.UnverifiedPayload{
		Action: DeleteUsersAction,
		Users:  unverifiedUsers,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal unverified users list")
	}

	err = p.sender.SendMessageToCustomChannel(p.unverifiedTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		return errors.Wrap(err, "failed to publish users", logan.F{
			"topic": p.unverifiedTopic,
		})
	}

	p.log.Infof("successfully published users to `%s`", p.unverifiedTopic)
	return nil
}

func (p *processor) buildMessage(uuid string, payload []byte) *message.Message {
	return &message.Message{
		UUID:     uuid,
		Metadata: nil,
		Payload:  payload,
	}
}

func (p *processor) sendUpdateUserSlack(uuid string, msg data.ModulePayload) error {
	marshaledPayload, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal update slack info")
	}

	err = p.sender.SendMessageToCustomChannel(p.identityTopic, p.buildMessage(uuid, marshaledPayload))
	if err != nil {
		return errors.Wrap(err, "failed to publish users", logan.F{
			"topic": p.identityTopic,
		})
	}

	p.log.Infof("successfully published user to `%s`", p.identityTopic)
	return nil
}
