/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Conversation struct {
	Key
	Attributes ConversationAttributes `json:"attributes"`
}
type ConversationResponse struct {
	Data     Conversation `json:"data"`
	Included Included     `json:"included"`
}

type ConversationListResponse struct {
	Data     []Conversation `json:"data"`
	Included Included       `json:"included"`
	Links    *Links         `json:"links"`
}

// MustConversation - returns Conversation from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustConversation(key Key) *Conversation {
	var conversation Conversation
	if c.tryFindEntry(key, &conversation) {
		return &conversation
	}
	return nil
}
