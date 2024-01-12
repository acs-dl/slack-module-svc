package requests

import (
	"encoding/json"
	"net/http"

	"github.com/acs-dl/slack-module-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type AddLinkRequest struct {
	Data resources.Link `json:"data"`
}

func NewAddLinkRequest(r *http.Request) (AddLinkRequest, error) {
	var request AddLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *AddLinkRequest) validate() error {
	return validation.Errors{
		"link": validation.Validate(&r.Data.Attributes.Link, validation.Required),
	}.Filter()
}
