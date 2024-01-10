package slack

import (
	"github.com/slack-go/slack"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *client) GetBillableInfo(priority int) (map[string]bool, error) {
	item, err := addFunctionInPQueue(c.pqueues.UserPQueue, c.userClient.GetBillableInfoForTeam, []any{}, priority)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	if err = item.Response.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get billable info from api")
	}

	billableInfo, ok := item.Response.Value.(map[string]slack.BillingActive)
	if !ok {
		return nil, errors.New("failed to convert response")
	}

	billableInfoResponse := make(map[string]bool)
	for user, bill := range billableInfo {
		billableInfoResponse[user] = bill.BillingActive
	}

	return billableInfoResponse, nil
}
