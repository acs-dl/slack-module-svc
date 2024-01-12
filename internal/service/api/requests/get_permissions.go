package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type GetPermissionsRequest struct {
	pgdb.OffsetPageParams

	Link     *string `filter:"link"`
	Username *string `filter:"username"`
}

func NewGetPermissionsRequest(r *http.Request) (GetPermissionsRequest, error) {
	var request GetPermissionsRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, errors.Wrap(err, "failed to decode url")
}
