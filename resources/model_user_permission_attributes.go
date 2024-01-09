/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserPermissionAttributes struct {
	AccessLevel AccessLevel `json:"access_level"`
	// is user billable
	Bill *bool `json:"bill,omitempty"`
	// chat title
	Link string `json:"link"`
	// user id from module
	ModuleId *string `json:"module_id,omitempty"`
	// workspace title
	Path string `json:"path"`
	// submodule id to handle submodule with the same title
	SubmoduleId *string `json:"submodule_id,omitempty"`
	// user id from identity
	UserId *int64 `json:"user_id,omitempty"`
	// username from slack
	Username *string `json:"username,omitempty"`
}
