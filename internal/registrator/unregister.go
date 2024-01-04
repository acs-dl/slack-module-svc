package registrator

import (
	"fmt"
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/data"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (r *registrar) UnregisterModule() error {
	r.logger.Infof("started unregister module `%s`", data.ModuleName)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", r.config.OuterUrl, data.ModuleName), nil)
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.From(errors.New("error in response"), logan.F{
			"status": res.Status,
		})
	}

	r.logger.Infof("finished unregister module `%s`", data.ModuleName)
	return nil
}
