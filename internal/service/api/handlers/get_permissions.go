package handlers

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/service/api/models"
	"github.com/acs-dl/slack-module-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPermissionsRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	permissionsQ := PermissionsQ(r).WithUsers()
	countPermissionsQ := PermissionsQ(r).CountWithUsers()

	if request.UserId != nil {
		permissionsQ = permissionsQ.FilterByUserIds(*request.UserId)
		countPermissionsQ = countPermissionsQ.FilterByUserIds(*request.UserId)
	}

	if request.Username != nil {
		permissionsQ = permissionsQ.FilterByUsernames(*request.Username)
		countPermissionsQ = countPermissionsQ.FilterByUsernames(*request.Username)
	}

	if request.ParentLink != nil {
		permissionsQ = permissionsQ.FilterByWorkspaceNames(*request.ParentLink)
		countPermissionsQ = countPermissionsQ.FilterByWorkspaceNames(*request.ParentLink)
	}

	if request.Link != nil {
		permissionsQ = permissionsQ.SearchBy(*request.Link)
		countPermissionsQ = countPermissionsQ.SearchBy(*request.Link)
	}

	permissions, err := permissionsQ.Page(request.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	totalCount, err := countPermissionsQ.GetTotalCount()
	if err != nil {
		Log(r).WithError(err).Error("failed to get permissions total count")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewUserPermissionListResponse(permissions)
	response.Meta.TotalCount = totalCount
	response.Links = data.GetOffsetLinksForPGParams(r, request.OffsetPageParams)

	ape.Render(w, response)
}
