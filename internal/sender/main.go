package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/slack-module-svc/internal/config"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const ServiceName = data.ModuleName + "-sender"

type Sender interface {
	Run(context.Context) error
	SendMessageToCustomChannel(topic string, msg *message.Message) error
}

type sender struct {
	publisher   *amqp.Publisher
	responsesQ  data.Responses
	log         *logan.Entry
	topic       string
	runnerDelay time.Duration
}

func New(cfg config.Config) Sender {
	return &sender{
		publisher:   cfg.Amqp().Publisher,
		responsesQ:  postgres.NewResponsesQ(cfg.DB()),
		log:         logan.New().WithField("service", ServiceName),
		topic:       cfg.Amqp().Orchestrator,
		runnerDelay: cfg.Runners().Sender,
	}
}

func (s *sender) Run(ctx context.Context) error {
	go running.WithBackOff(ctx, s.log,
		ServiceName,
		s.processMessages,
		s.runnerDelay,
		s.runnerDelay,
		s.runnerDelay,
	)

	return nil
}

func (s *sender) processMessages(ctx context.Context) error {
	s.log.Info("started processing responses")

	responses, err := s.responsesQ.Select()
	if err != nil {
		return errors.Wrap(err, "failed to select responses")
	}

	for _, response := range responses {
		s.log.Infof("started processing response with id %s", response.ID)
		err = (*s.publisher).Publish(s.topic, s.buildResponse(response))
		if err != nil {
			return errors.Wrap(err, "failed to process response", logan.F{
				"response_id": response.ID,
			})
		}

		err = s.responsesQ.FilterByIds(response.ID).Delete()
		if err != nil {
			return errors.Wrap(err, "failed to delete processed response `%s", logan.F{
				"response_id": response.ID,
			})
		}
		s.log.Info("finished processing response with id ", response.ID)
	}

	s.log.Info("finished processing responses")
	return nil
}

func (s *sender) buildResponse(response data.Response) *message.Message {
	marshaled, err := json.Marshal(response)
	if err != nil {
		s.log.WithError(err).Errorf("failed to marshal response")
	}

	return &message.Message{
		UUID:     response.ID,
		Metadata: nil,
		Payload:  marshaled,
	}
}

func (s *sender) SendMessageToCustomChannel(topic string, msg *message.Message) error {
	err := (*s.publisher).Publish(topic, msg)
	if err != nil {
		return errors.Wrap(err, "failed to send msg `%s to `%s`", logan.F{
			"message_id": msg.UUID,
			"topic":      topic,
		})
	}

	return nil
}
