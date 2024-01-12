package worker

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (w *worker) validateRefreshSubmodulesRequest(msg data.ModulePayload) error {
	return validation.Errors{
		"links": validation.Validate(msg.Links, validation.Required),
	}.Filter()
}
