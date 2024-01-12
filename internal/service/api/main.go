package api

import (
	"context"
	"net/http"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Router interface {
	Run(ctx context.Context) error
}

type router struct {
	cfg config.Config
	ctx context.Context
}

func (r *router) Run(_ context.Context) error {
	router := r.apiRouter()

	if err := r.cfg.Copus().RegisterChi(router); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(r.cfg.Listener(), router)
}

func New(cfg config.Config, ctx context.Context) Router {
	return &router{
		cfg: cfg,
		ctx: ctx,
	}
}
