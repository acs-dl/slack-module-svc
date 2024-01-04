package worker

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (w *worker) sendUsers(uuid string, users []data.User) error {
	unverifiedUsers := make([]data.UnverifiedUser, 0)
	for i := range users {
		if users[i].Id != nil {
			continue
		}

		permission, err := w.permissionsQ.FilterBySlackIds(users[i].SlackId).FilterByGreaterTime(users[i].CreatedAt).Get()
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

	err = w.sender.SendMessageToCustomChannel(w.unverifiedTopic, w.buildMessage(uuid, marshaledPayload))
	if err != nil {
		return errors.Wrap(err, "failed to publish users", logan.F{
			"topic": w.unverifiedTopic,
		})
	}

	w.logger.Infof("successfully published users to `unverified-svc`")
	return nil
}

func (w *worker) buildMessage(uuid string, payload []byte) *message.Message {
	return &message.Message{
		UUID:     uuid,
		Metadata: nil,
		Payload:  payload,
	}
}
