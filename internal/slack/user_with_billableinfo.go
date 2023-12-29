package slack

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (s *client) BillableInfoForUser(userId string) (bool, error) {
	userBill, err := s.userAdminClient.GetBillableInfo(userId)
	if err != nil {
		return false, errors.Wrap(err, "failed to get billable info about user", logan.F{
			"user_id": userId,
		})
	}

	return userBill[userId].BillingActive, nil
}
