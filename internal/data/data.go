package data

const (
	ModuleName        = "slack"
	UnverifiedService = "unverified-svc"
	IdentityService   = "identity"
)

type ModulePayload struct {
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
	Action    string `json:"action"`

	//other fields that are required for module
	Link        string   `json:"link"`
	SlackId     string   `json:"slack_id"`
	Username    string   `json:"username"`
	Realname    *string  `json:"real_name"`
	AccessLevel string   `json:"access_level"`
	Links       []string `json:"links"`
}

type UnverifiedPayload struct {
	Action string           `json:"action"`
	Users  []UnverifiedUser `json:"users"`
}

var Roles = map[string]string{
	"":             "No access",
	Admin:          "Admin",
	Member:         "User",
	Owner:          "Owner",
	RestrictedUser: "Restricted User",
	Bot:            "Bot",
	App:            "App",
}

func MapKeysToSlice(m map[string]string) []interface{} {
	keys := make([]interface{}, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
