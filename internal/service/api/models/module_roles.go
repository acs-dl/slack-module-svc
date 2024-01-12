package models

import "github.com/acs-dl/slack-module-svc/resources"

func NewModuleRolesResponse() resources.ModuleRolesResponse {
	return resources.ModuleRolesResponse{
		Data: resources.ModuleRoles{
			Key:        resources.NewKeyInt64(0, resources.MODULES),
			Attributes: map[string]interface{}{},
		},
	}
}
