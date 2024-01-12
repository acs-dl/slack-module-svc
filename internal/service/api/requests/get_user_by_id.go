package requests

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewGetUserByIdRequest(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")

	return validateUserId(id)
}

func validateUserId(id string) (int64, error) {
	err := validation.Errors{
		"id": validation.Validate(id, validation.Required),
	}.Filter()

	if err != nil {
		return 0, errors.Wrap(err, "id not provided", logan.F{
			"id": id,
		})
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)
	return parsedId, errors.Wrap(err, "failed to parse int", logan.F{
		"id": id,
	})
}
