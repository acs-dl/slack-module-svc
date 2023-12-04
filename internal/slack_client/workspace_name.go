package slack_client

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) WorkspaceName() (string, error) {

	team, err := s.userAdminClient.GetTeamInfo()
	if err != nil {
		return "", errors.Wrap(err, "Error retrieving team info for the transferred token.")
	}

	return team.Name, nil
}
