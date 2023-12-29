package handlers

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/service/api/models"
	"github.com/acs-dl/slack-module-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRole(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRoleRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, models.NewRoleResponse(data.Roles[*request.AccessLevel], *request.AccessLevel))
}
