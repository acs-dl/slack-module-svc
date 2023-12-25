package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Users interface {
	New() Users

	Upsert(user User) error
	Delete() error
	Select() ([]User, error)
	Get() (*User, error)

	FilterByLowerTime(time time.Time) Users
	FilterById(id *int64) Users
	FilterBySlackIds(telegramIds ...string) Users
	FilterByUsername(username string) Users
	//FilterByPhone(phone string) Users
	SearchBy(search string) Users

	Count() Users
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) Users
}

type User struct {
	Id        *int64    `json:"-" db:"id" structs:"id,omitempty"`
	Username  *string   `json:"username" db:"username" structs:"username,omitempty"`    //name from slack_client
	Realname  *string   `json:"real_name" db:"real_name" structs:"real_name,omitempty"` //real_name from slack_client
	SlackId   string    `json:"slack_id" db:"slack_id" structs:"slack_id,omitempty"`    //id from slack_client
	CreatedAt time.Time `json:"created_at" db:"created_at" structs:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" structs:"-"`
	// fields to create permission
	AccessLevel string `json:"-" db:"-" structs:"-"`
}

type UnverifiedUser struct {
	CreatedAt time.Time `json:"created_at"`
	Module    string    `json:"module"`
	Submodule string    `json:"submodule"`
	Username  *string   `json:"username,omitempty"`
	RealName  *string   `json:"real_name,omitempty"`
	SlackId   string    `json:"slack_id"`
	ModuleId  string    `json:"module_id"`
}
