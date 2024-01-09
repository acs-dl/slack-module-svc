package models

import "github.com/acs-dl/slack-module-svc/resources"

type ModuleRolesResponse struct {
	Data ModuleRoles `json:"data"`
}

type Roles map[string]string
type ModuleRoles struct {
	resources.Key
	Attributes Roles `json:"attributes"`
}

func NewModuleRolesResponse() ModuleRolesResponse {
	return ModuleRolesResponse{
		Data: ModuleRoles{
			Key: resources.Key{
				ID:   "0",
				Type: resources.MODULES,
			},
			Attributes: Roles{},
		},
	}
}
