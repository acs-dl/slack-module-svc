package processor

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) getUserFromDbBySlackId(slackId string) (*data.User, error) {
	user, err := p.managerQ.Users.FilterBySlackIds(slackId).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db", logan.F{
			"slack_id": slackId,
		})
	}

	if user == nil {
		return nil, errors.From(errors.New("user not found in module"), logan.F{
			"slack_id": slackId,
		})
	}

	return user, nil
}
