package slack_client

import (
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *slackStruct) BillableInfoForUser(userId string) (bool, error) {
	userBill, err := s.userAdminClient.GetBillableInfo(userId)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("failed to get billable info about user %s.", userId))
	}

	return userBill[userId].BillingActive, nil
}
