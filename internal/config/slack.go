package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const slackYamlTag = "slack"

type SlackConfig struct {
	UserToken string `fig:"user_oauth_token,required"`
	BotToken  string `fig:"bot_user_oauth_token,required"`
}

func (c *config) SlackParams() SlackConfig {
	return c.slackCfg.Do(func() interface{} {
		cfg := SlackConfig{}

		if err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, slackYamlTag)).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out slack tokens"))
		}

		return cfg
	}).(SlackConfig)
}
