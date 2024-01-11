package processor

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (p *processor) validateVerifyUser(msg data.ModulePayload) error {
	return validation.Errors{
		"user_id":  validation.Validate(msg.UserId, validation.Required),
		"username": validation.Validate(msg.Username, validation.Required),
	}.Filter()
}
