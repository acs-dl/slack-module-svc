package slack

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetWorkspaceName() (string, error) {
	team, err := s.botClient.GetTeamInfo()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve team info")
	}

	return team.Name, nil
}
