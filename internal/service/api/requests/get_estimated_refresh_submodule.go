package requests

import (
	"encoding/json"
	"net/http"

	"github.com/acs-dl/slack-module-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type RefreshSubmoduleRequest struct {
	Data resources.Submodules `json:"data"`
}

func NewRefreshSubmoduleRequest(r *http.Request) (RefreshSubmoduleRequest, error) {
	var request RefreshSubmoduleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *RefreshSubmoduleRequest) validate() error {
	return validation.Errors{
		"links": validation.Validate(&r.Data.Attributes.Links, validation.Required),
	}.Filter()
}
