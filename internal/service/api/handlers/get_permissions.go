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

	var usernames []string
	if request.Username != nil {
		usernames = append(usernames, *request.Username)
	}

	var ids []int64
	if request.UserId != nil {
		ids = append(ids, *request.UserId)
	}

	permissionsQ := PermissionsQ(r).WithUsers().FilterByUsernames(usernames...).FilterByUserIds(ids...)
	countPermissionsQ := PermissionsQ(r).CountWithUsers().FilterByUsernames(usernames...).FilterByUserIds(ids...)

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
