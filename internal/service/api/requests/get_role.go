package requests

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval"
)

type GetRoleRequest struct {
	AccessLevel *string `filter:"accessLevel"`
}

func NewGetRoleRequest(r *http.Request) (GetRoleRequest, error) {
	var request GetRoleRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, validateGetRoleRequest(request)
}

func validateGetRoleRequest(request GetRoleRequest) error {
	return validation.Errors{
		"accessLevel": validation.Validate(request.AccessLevel, validation.Required),
	}.Filter()
}
