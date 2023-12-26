package slack_client

import "github.com/slack-go/slack"

func (s *slackStruct) GetBillableInfo() (map[string]slack.BillingActive, error) {
	return s.userAdminClient.GetBillableInfoForTeam()
}
