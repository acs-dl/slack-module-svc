package worker

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (w *worker) getConversationsByLink(link string) ([]data.Conversation, error) {
	return helpers.GetConversationsByLink(
		w.pqueues.BotPQueue,
		any(w.client.GetConversationsByLink),
		[]interface{}{link},
		pqueue.LowPriority,
	)
}

func (w *worker) validateRefreshSubmodulesRequest(msg data.ModulePayload) error {
	return validation.Errors{
		"links": validation.Validate(msg.Links, validation.Required),
	}.Filter()
}
