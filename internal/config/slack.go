package config

import (
	"os"
)

const slackYamlTag = "slack"

type SlackConfig struct {
	UserToken string
	BotToken  string
}

func (c *config) SlackParams() SlackConfig {
	return c.slackCfg.Do(func() interface{} {
		userToken := os.Getenv("USER_TOKEN")
		if userToken == "" {
			panic("no user token was provided")
		}

		botToken := os.Getenv("BOT_TOKEN")
		if botToken == "" {
			panic("no bot token was provided")
		}

		return SlackConfig{
			UserToken: userToken,
			BotToken:  botToken,
		}
	}).(SlackConfig)
}
