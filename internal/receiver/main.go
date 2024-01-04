package receiver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
	"github.com/acs-dl/slack-module-svc/internal/processor"
	"github.com/acs-dl/slack-module-svc/internal/worker"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	ServiceName = data.ModuleName + "-receiver"

	VerifyUserAction       = "verify_user"
	RefreshModuleAction    = "refresh_module"
	RefreshSubmoduleAction = "refresh_submodule"
	DeleteUserAction       = "delete_user"
)

type Receiver interface {
	HandleNewMessage(msg data.ModulePayload) (string, error)
	Run(ctx context.Context) error
}

type receiver struct {
	subscriber  *amqp.Subscriber
	topic       string
	log         *logan.Entry
	processor   processor.Processor
	worker      worker.Worker
	responseQ   data.Responses
	runnerDelay time.Duration
}

var handleActions = map[string]func(r *receiver, msg data.ModulePayload) (string, error){
	VerifyUserAction: func(r *receiver, msg data.ModulePayload) (string, error) {
		return r.processor.HandleVerifyUserAction(msg)
	},
	RefreshModuleAction: func(r *receiver, msg data.ModulePayload) (string, error) {
		return r.worker.RefreshModule()
	},
	RefreshSubmoduleAction: func(r *receiver, msg data.ModulePayload) (string, error) {
		return r.worker.RefreshSubmodules(msg)
	},
}

func New(cfg config.Config, ctx context.Context) Receiver {
	return &receiver{
		subscriber:  cfg.Amqp().Subscriber,
		topic:       cfg.Amqp().Topic,
		log:         logan.New().WithField("service", ServiceName),
		processor:   processor.ProcessorInstance(ctx),
		responseQ:   postgres.NewResponsesQ(cfg.DB()),
		worker:      worker.WorkerInstance(ctx),
		runnerDelay: cfg.Runners().Receiver,
	}
}

func (r *receiver) Run(ctx context.Context) error {
	go running.WithBackOff(ctx, r.log,
		ServiceName,
		r.listenMessages,
		r.runnerDelay,
		r.runnerDelay,
		r.runnerDelay,
	)

	return nil
}

func validateMessageAction(msg data.ModulePayload) error {
	return validation.Errors{
		"action": validation.Validate(msg.Action, validation.Required, validation.In(VerifyUserAction, RefreshModuleAction, RefreshSubmoduleAction)),
	}.Filter()
}

func (r *receiver) listenMessages(ctx context.Context) error {
	r.log.Info("started listening messages")
	return r.subscribeForTopic(ctx, r.topic)
}

func (r *receiver) subscribeForTopic(ctx context.Context, topic string) error {
	msgChan, err := r.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe for topic", logan.F{
			"topic": topic,
		})
	}
	r.log.Info("successfully subscribed for topic ", topic)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgChan:
			r.log.Info("received message ", msg.UUID)
			err = r.processMessage(msg)
			if err != nil {
				r.log.WithError(err).Error("failed to process message ", msg.UUID)
			}
			msg.Ack()
		}
	}
}

func (r *receiver) HandleNewMessage(msg data.ModulePayload) (string, error) {
	r.log.Infof("handling message with id `%s`", msg.RequestId)

	err := validateMessageAction(msg)
	if err != nil {
		return data.FAILURE, errors.Wrap(err, "no such action to handle for message", logan.F{
			"action": msg.RequestId,
		})
	}

	requestHandler := handleActions[msg.Action]
	requestStatus, err := requestHandler(r, msg)
	if err != nil {
		return requestStatus, errors.Wrap(err, "failed to handle message", logan.F{
			"action": msg.RequestId,
		})
	}

	r.log.Infof("finish handling message with id `%s`", msg.RequestId)
	return requestStatus, nil
}

func (r *receiver) processMessage(msg *message.Message) error {
	r.log.Info("started processing message ", msg.UUID)

	var queueOutput data.ModulePayload
	err := json.Unmarshal(msg.Payload, &queueOutput)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message", logan.F{
			"message_uuid": msg.UUID,
		})
	}

	queueOutput.RequestId = msg.UUID
	var errMsg *string = nil
	requestStatus, err := r.HandleNewMessage(queueOutput)
	if err != nil {
		requestError := err.Error()
		errMsg = &requestError
		r.log.WithError(err).Error("failed to process message ", msg.UUID)
	}

	err = r.responseQ.Insert(data.Response{
		ID:      msg.UUID,
		Status:  requestStatus,
		Error:   errMsg,
		Payload: json.RawMessage(msg.Payload),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create response", logan.F{
			"message_uuid": msg.UUID,
		})
	}

	r.log.Info("finished processing message ", msg.UUID)
	return nil
}
