package registrator

import (
	"context"
	"time"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/distributed_lab/running"
)

const ServiceName = data.ModuleName + "-registrar"

type Registrar interface {
	Run(ctx context.Context) error
	UnregisterModule() error
}

type registrar struct {
	logger      *logan.Entry
	config      config.RegistratorConfig
	runnerDelay time.Duration
}

func New(cfg config.Config) Registrar {
	return &registrar{
		logger:      cfg.Log().WithField("runner", ServiceName),
		config:      cfg.Registrator(),
		runnerDelay: cfg.Runners().Registrar,
	}
}

func (r *registrar) Run(ctx context.Context) error {
	running.WithBackOff(
		ctx,
		r.logger,
		ServiceName,
		r.registerModule,
		r.runnerDelay,
		r.runnerDelay,
		r.runnerDelay,
	)

	return nil
}
