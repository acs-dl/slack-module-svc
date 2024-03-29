package data

const (
	ModuleName = "slack"
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

func GetRoles() []interface{} {
	roles := make([]interface{}, 0, len(Roles))
	for role := range Roles {
		roles = append(roles, role)
	}
	return roles
}
