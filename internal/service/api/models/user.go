package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func NewUserModel(user data.User, id int) resources.User {
	userAccessLevel := data.Roles[user.AccessLevel]
	result := resources.User{
		Key: resources.NewKeyInt64(int64(id), resources.USER),
		Attributes: resources.UserAttributes{
			UserId:      user.Id,
			Username:    *user.Username,
			SlackId:     &user.SlackId,
			Module:      data.ModuleName,
			CreatedAt:   &user.CreatedAt,
			AccessLevel: &userAccessLevel,
		},
	}

	return result
}

func NewUserResponse(user data.User) resources.UserResponse {
	return resources.UserResponse{
		Data: NewUserModel(user, 0),
	}
}

type Meta struct {
	TotalCount int64 `json:"total_count"`
}
