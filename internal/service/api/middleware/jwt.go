package middleware

import (
	"net/http"

	auth "github.com/acs-dl/auth-svc/middlewares"
	"github.com/acs-dl/slack-module-svc/internal/data"
)

func IsAdmin(secret string) func(http.Handler) http.Handler {
	return auth.Jwt(secret, data.ModuleName, []string{data.Roles[data.Admin], data.Roles[data.Owner]}...)
}

func IsAuthenticated(secret string) func(http.Handler) http.Handler {
	return auth.Jwt(secret, data.ModuleName, []string{data.Roles[data.Admin], data.Roles[data.Owner], data.Roles[data.Member]}...)
}
