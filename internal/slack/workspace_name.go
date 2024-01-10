package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetWorkspaceName(priority int) (string, error) {
	item, err := addFunctionInPQueue(c.pqueues.BotPQueue, c.botClient.GetTeamInfo, []any{}, priority)
	if err != nil {
		return "", errors.Wrap(err, "failed to add function in pqueue")
	}

	if err = item.Response.Error; err != nil {
		return "", errors.Wrap(err, "failed to get team info from api")
	}

	teamInfo, ok := item.Response.Value.(*slack.TeamInfo)
	if !ok {
		return "", errors.New("failed to convert response")
	}

	return teamInfo.Name, nil
}
