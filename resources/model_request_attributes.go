/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RequestAttributes struct {
	// Module to grant permission
	Module string `json:"module"`
	// Already built payload to grant permission <br><br> -> \"get_users\" = action to get users with their permissions from channel in slack<br> -> \"verify_user\" = action to verify user in slack module (connect user id from identity with slack info)<br> -> \"delete_user\" = action to delete user from module (from all links)<br>
	Payload json.RawMessage `json:"payload"`
}
