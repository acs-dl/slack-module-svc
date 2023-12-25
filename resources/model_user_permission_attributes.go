/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserPermissionAttributes struct {
	AccessLevel AccessLevel `json:"access_level"`
	// chat title
	Link string `json:"link"`
	// user id from module
	ModuleId *int64 `json:"module_id,omitempty"`
	// workspace title
	Path string `json:"path"`
	// id from slack
	SlackId *string `json:"slack_id,omitempty"`
	// submodule id to handle submodule with the same title
	SubmoduleId *string `json:"submodule_id,omitempty"`
	// user id from identity
	UserId *int64 `json:"user_id,omitempty"`
	// username from slack
	Username *string `json:"username,omitempty"`
}
