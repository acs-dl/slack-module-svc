package handlers

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/ape"
)

func GetRolesMap(w http.ResponseWriter, r *http.Request) {
	result := newModuleRolesResponse()

	for key, val := range data.Roles {
		result.Data.Attributes[key] = val
	}

	ape.Render(w, result)
}
