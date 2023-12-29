package slack

import "github.com/slack-go/slack"

func (s *client) GetBillableInfo() (map[string]slack.BillingActive, error) {
	return s.userAdminClient.GetBillableInfoForTeam()
}
