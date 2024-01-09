package requests

import (
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type GetRoleRequest struct {
	AccessLevel *string `filter:"accessLevel"`
}

func NewGetRoleRequest(r *http.Request) (GetRoleRequest, error) {
	var request GetRoleRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, errors.Wrap(err, "failed to decode url")
	}

	return request, request.validate()
}

func (r *GetRoleRequest) validate() error {
	return validation.Errors{
		"accessLevel": validation.Validate(r.AccessLevel, validation.Required, validation.In(data.GetRoles()...)),
	}.Filter()
}
