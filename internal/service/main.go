package service

import (
	"context"
	"sync"

	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"github.com/acs-dl/slack-module-svc/internal/processor"
	"github.com/acs-dl/slack-module-svc/internal/receiver"
	"github.com/acs-dl/slack-module-svc/internal/registrator"
	"github.com/acs-dl/slack-module-svc/internal/sender"
	"github.com/acs-dl/slack-module-svc/internal/service/api"
	"github.com/acs-dl/slack-module-svc/internal/service/api/handlers"
	"github.com/acs-dl/slack-module-svc/internal/service/types"
	"github.com/acs-dl/slack-module-svc/internal/worker"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type runner struct {
	name    string
	service types.Service
}

func Run(cfg config.Config) {
	logger := cfg.Log().WithField("service", "main")
	ctx := context.Background()
	wg := new(sync.WaitGroup)

	logger.Info("Starting all available services...")

	stopProcessQueue := make(chan struct{})
	pqueues := pqueue.NewPQueues(logger)
	go pqueues.BotPQueue.ProcessQueue(cfg.RateLimit().RequestsAmount, cfg.RateLimit().TimeLimit, stopProcessQueue)
	go pqueues.UserPQueue.ProcessQueue(cfg.RateLimit().RequestsAmount, cfg.RateLimit().TimeLimit, stopProcessQueue)
	ctx = pqueue.CtxPQueues(&pqueues, ctx)
	ctx = handlers.CtxConfig(cfg, ctx)

	// Services instantiation goes below
	senderSvc := sender.New(cfg)
	ctx = sender.CtxSenderInstance(senderSvc, ctx)

	processorSvc := processor.New(cfg, ctx)
	ctx = processor.CtxProcessorInstance(processorSvc, ctx)

	workerSvc := worker.New(cfg, ctx)
	ctx = worker.CtxWorkerInstance(workerSvc, ctx)

	receiverSvc := receiver.New(cfg, ctx)
	ctx = receiver.CtxReceiverInstance(receiverSvc, ctx)

	runners := []runner{
		{"sender", senderSvc},
		{"worker", workerSvc},
		{"receiver", receiverSvc},
		{"registrar", registrator.New(cfg)},
		{"api", api.New(cfg, ctx)},
	}

	for _, runner := range runners {
		wg.Add(1)
		go func(svc types.Service, ctx context.Context) {
			defer wg.Done()

			if err := svc.Run(ctx); err != nil {
				panic(errors.Wrap(err, "failed to run service", logan.F{
					"serviceName": runner.name,
				}))
			}
		}(runner.service, ctx)

		logger.WithField("service", runner.name).Info("Service started")
	}

	wg.Wait()
}
