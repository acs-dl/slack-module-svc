package processor

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (p *processor) getUser(slackID string) (*data.User, error) {
	user, err := p.client.GetUser(slackID, pqueue.NormalPriority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from api", logan.F{
			"slack_id": slackID,
		})
	}
	return user, nil
}

func (p *processor) validateVerifyUser(msg data.ModulePayload) error {
	return validation.Errors{
		"user_id":  validation.Validate(msg.UserId, validation.Required),
		"username": validation.Validate(msg.Username, validation.Required),
	}.Filter()
}
