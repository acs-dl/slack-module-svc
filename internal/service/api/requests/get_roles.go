package requests

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type GetRolesRequest struct {
	Link     *string `filter:"link"`
	Username *string `filter:"username"`
}

func NewGetRolesRequest(r *http.Request) (GetRolesRequest, error) {
	var request GetRolesRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, errors.Wrap(err, "failed to decode url")
	}

	return request, request.validate()
}

func (r *GetRolesRequest) validate() error {
	return validation.Errors{
		"link":     validation.Validate(&r.Link, validation.Required),
		"username": validation.Validate(&r.Username, validation.Required),
	}.Filter()
}
