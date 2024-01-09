package models

import "github.com/acs-dl/slack-module-svc/resources"

func NewRolesModel(found bool, roles []resources.AccessLevel) resources.Roles {
	result := resources.Roles{
		Key: resources.Key{
			ID:   "roles",
			Type: resources.ROLES,
		},
		Attributes: resources.RolesAttributes{
			Req:  found,
			List: roles,
		},
	}

	return result
}
