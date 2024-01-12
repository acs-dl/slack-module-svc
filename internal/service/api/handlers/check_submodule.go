package handlers

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/service/api/models"
	"github.com/acs-dl/slack-module-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CheckSubmodule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCheckSubmoduleRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	conversations, err := ConversationsQ(r).SearchBy(*request.Link).Select()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get conversations with `%s` title", *request.Link)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if len(conversations) == 0 {
		Log(r).Warnf("no conversations were found")
		ape.Render(w, models.NewLinkResponse("", false, conversations))
		return
	}

	ape.Render(w, models.NewLinkResponse(*request.Link, true, conversations))
}
