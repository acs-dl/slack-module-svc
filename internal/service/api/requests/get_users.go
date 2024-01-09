package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type GetUsersRequest struct {
	Username *string `filter:"username"`
}

func NewGetUsersRequest(r *http.Request) (GetUsersRequest, error) {
	var request GetUsersRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, errors.Wrap(err, "failed to decode url")
}
