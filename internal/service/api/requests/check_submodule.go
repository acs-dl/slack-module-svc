package requests

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type CheckSubmoduleRequest struct {
	Link *string `filter:"link"`
}

func NewCheckSubmoduleRequest(r *http.Request) (CheckSubmoduleRequest, error) {
	var request CheckSubmoduleRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, errors.Wrap(err, "failed to decode url")
	}

	return request, request.validate()
}

func (r *CheckSubmoduleRequest) validate() error {
	return validation.Errors{
		"link": validation.Validate(&r.Link, validation.Required),
	}.Filter()
}
