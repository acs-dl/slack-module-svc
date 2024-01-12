package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func NewUserInfoModel(user data.User, id int) resources.UserInfo {
	return resources.UserInfo{
		Key: resources.NewKeyInt64(int64(id), resources.USER),
		Attributes: resources.UserInfoAttributes{
			Username: *user.Username,
			SlackId:  user.SlackId,
			Name:     *user.Realname,
		},
	}
}

func NewUserInfoList(users []data.User, offset uint64) []resources.UserInfo {
	result := make([]resources.UserInfo, len(users))
	for i, user := range users {
		result[i] = NewUserInfoModel(user, i+int(offset))
	}

	return result
}

func NewUserInfoListResponse(users []data.User, offset uint64) resources.UserInfoListResponse {
	return resources.UserInfoListResponse{
		Data: NewUserInfoList(users, offset),
	}
}
