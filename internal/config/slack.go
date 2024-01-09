package config

import (
	"os"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

const (
	userTokenRegex = "xoxp-[0-9]{13}-[0-9]{13}-[0-9]{13}-[0-9a-zA-Z]{32}"
	botTokenRegex  = "xoxb-[0-9]{13}-[0-9]{13}-[0-9a-zA-Z]{24}"
)

type SlackConfig struct {
	UserToken string
	BotToken  string
}

func (c *config) SlackParams() SlackConfig {
	return c.slackCfg.Do(func() interface{} {
		config := SlackConfig{
			UserToken: os.Getenv("USER_TOKEN"),
			BotToken:  os.Getenv("BOT_TOKEN"),
		}

		if err := config.validate(); err != nil {
			panic(errors.Wrap(err, "slack oauth tokens validation failed"))
		}

		return config
	}).(SlackConfig)
}

func (cfg *SlackConfig) validate() error {
	return validation.Errors{
		"USER_TOKEN": validation.Validate(&cfg.UserToken, validation.Required, validation.Match(regexp.MustCompile(userTokenRegex))),
		"BOT_TOKEN":  validation.Validate(&cfg.BotToken, validation.Required, validation.Match(regexp.MustCompile(botTokenRegex))),
	}.Filter()
}
