package worker

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (w *worker) getUserFromDbBySlackId(slackId string) (*data.User, error) {
	usersQ := w.usersQ.New()
	user, err := usersQ.FilterBySlackIds(slackId).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db")
	}

	if user == nil {
		return nil, errors.Errorf("no such user in module")
	}

	return user, nil
}
