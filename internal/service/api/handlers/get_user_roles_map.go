package handlers

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/service/api/models"
	"gitlab.com/distributed_lab/ape"
)

func GetUserRolesMap(w http.ResponseWriter, r *http.Request) {
	result := models.NewModuleRolesResponse()

	result.Data.Attributes["super_admin"] = data.Roles[data.Owner]
	result.Data.Attributes["admin"] = data.Roles[data.Admin]
	result.Data.Attributes["user"] = data.Roles[data.Member]

	ape.Render(w, result)
}
