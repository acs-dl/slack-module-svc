package handlers

import (
	"net/http"
	"strings"

	"github.com/acs-dl/slack-module-svc/internal/service/api/models"
	"github.com/acs-dl/slack-module-svc/internal/service/api/requests"
	"github.com/acs-dl/slack-module-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRolesRequest(r)
	if err != nil {
		Log(r).WithError(err).Errorf("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.Link == nil {
		Log(r).Warnf("no link was provided")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	link := strings.ToLower(*request.Link)

	if request.Username == nil {
		Log(r).Warnf("no username was provided")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	permission, err := PermissionsQ(r).FilterByUsernames(*request.Username).FilterByLinks(link).Get()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get permission from `%s` to `%s`", link, *request.Username)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if permission == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, models.NewRolesModel(true, []resources.AccessLevel{}))
}
