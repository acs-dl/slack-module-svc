package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetWorkspaceName(priority int) (string, error) {
	var teamInfo *slack.TeamInfo
	err := doQueueRequest[*slack.TeamInfo](QueueParameters{
		queue:    c.pqueues.BotPQueue,
		function: c.botClient.GetTeamInfo,
		args:     []any{},
		priority: priority,
	}, &teamInfo)
	if err != nil {
		return "", errors.Wrap(err, "failed to get team info")
	}

	return teamInfo.Name, nil
}
