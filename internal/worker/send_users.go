package worker

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (w *Worker) sendUsers(uuid string, users []data.User) error {
	unverifiedUsers := make([]data.UnverifiedUser, 0)
	for i := range users {
		if users[i].Id != nil {
			continue
		}
		permission, err := w.permissionsQ.FilterBySlackIds(users[i].SlackId).FilterByGreaterTime(users[i].CreatedAt).Get()

		if err != nil {
			w.logger.WithError(err).Errorf("failed to select permissions by date `%s`", users[i].CreatedAt.String())
			return errors.Wrap(err, "failed to select permissions by date")
		}

		if permission == nil {
			continue
		}

		unverifiedUsers = append(unverifiedUsers, convertUserToUnverifiedUser(users[i], permission.Link))
	}

	marshaledPayload, err := json.Marshal(data.UnverifiedPayload{
		Action: SetUsersAction,
		Users:  unverifiedUsers,
	})
	if err != nil {
		w.logger.WithError(err).Errorf("failed to marshal unverified users list")
		return errors.Wrap(err, "failed to marshal unverified users list")
	}

	err = w.sender.SendMessageToCustomChannel(data.UnverifiedService, w.buildMessage(uuid, marshaledPayload))
	if err != nil {
		w.logger.WithError(err).Errorf("failed to publish users to `slack-module`")
		return errors.Wrap(err, "failed to publish users to `slack-module`")
	}

	w.logger.Infof("successfully published users to `unverified-svc`")
	return nil
}

func (w *Worker) buildMessage(uuid string, payload []byte) *message.Message {
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
