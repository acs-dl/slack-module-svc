package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetBillableInfo(priority int) (map[string]bool, error) {
	var billableInfo map[string]slack.BillingActive
	err := doQueueRequest[map[string]slack.BillingActive](QueueParameters{
		queue:    c.pqueues.UserPQueue,
		function: c.userClient.GetBillableInfoForTeam,
		args:     []any{},
		priority: priority,
	}, &billableInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get billable info")
	}

	billableInfoResponse := make(map[string]bool)
	for user, bill := range billableInfo {
		billableInfoResponse[user] = bill.BillingActive
	}

	return billableInfoResponse, nil
}
