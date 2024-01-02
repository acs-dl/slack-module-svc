package slack

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) GetBillableInfo() (map[string]bool, error) {
	billableInfo, err := s.userClient.GetBillableInfoForTeam()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get billable info for team")
	}

	billableInfoResponse := make(map[string]bool)
	for user, bill := range billableInfo {
		billableInfoResponse[user] = bill.BillingActive
	}

	return billableInfoResponse, nil
}
