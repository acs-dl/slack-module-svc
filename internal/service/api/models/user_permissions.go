package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func NewUserPermissionModel(permission data.Permission, counter int) resources.UserPermission {
	result := resources.UserPermission{
		Key: resources.NewKeyInt64(int64(counter), resources.USER_PERMISSION),
		Attributes: resources.UserPermissionAttributes{
			AccessLevel: resources.AccessLevel{
				Name:  data.Roles[permission.AccessLevel],
				Value: permission.AccessLevel,
			},
			Link:        permission.Link,
			Username:    &permission.Username,
			ModuleId:    &permission.SlackId,
			UserId:      permission.Id,
			SubmoduleId: &permission.SubmoduleId,
			Path:        permission.Link,
			Bill:        &permission.Bill,
		},
	}

	return result
}

func NewUserPermissionList(permissions []data.Permission) []resources.UserPermission {
	result := make([]resources.UserPermission, len(permissions))
	for i, permission := range permissions {
		result[i] = NewUserPermissionModel(permission, i)
	}

	return result
}

func NewUserPermissionListResponse(permissions []data.Permission) UserPermissionListResponse {
	return UserPermissionListResponse{
		Data: NewUserPermissionList(permissions),
	}
}

type UserPermissionListResponse struct {
	Meta  Meta                       `json:"meta"`
	Data  []resources.UserPermission `json:"data"`
	Links *resources.Links           `json:"links"`
}
