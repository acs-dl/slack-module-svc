package data

type Conversations interface {
	New() Conversations

	Upsert(chat Conversation) error
	Delete() error
	Get() (*Conversation, error)
	Select() ([]Conversation, error)
	SearchBy(search string) Conversations

	FilterByTitles(titles ...string) Conversations
	FilterByIds(ids ...string) Conversations
}

type Conversation struct {
	Title         string `json:"title" db:"title" structs:"title"`
	Id            string `json:"id" db:"id" structs:"id"`
	MembersAmount int64  `json:"members_amount" db:"members_amount" structs:"members_amount"`
}
