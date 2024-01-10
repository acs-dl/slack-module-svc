package slack

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetWorkspaceName() (string, error) {
	team, err := c.botClient.GetTeamInfo()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve team info")
	}

	return team.Name, nil
}
