/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserPermissionAttributes struct {
	AccessLevel AccessLevel `json:"access_level"`
	// chat title
	Link string `json:"link"`
	// workspace title
	Path string `json:"path"`
	// user id from identity
	UserId *int64 `json:"user_id,omitempty"`
	// username from telegram
	Username *string `json:"username,omitempty"`
}
