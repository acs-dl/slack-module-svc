package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	Owner          = "owner"
	Admin          = "admin"
	Member         = "member"
	RestrictedUser = "restricted_user"
	Bot            = "bot"
	App            = "app"
)

type Permissions interface {
	New() Permissions

	Upsert(permission Permission) error
	UpdateAccessLevel(permission Permission) error
	Delete() error
	Select() ([]Permission, error)
	Get() (*Permission, error)

	FilterBySlackIds(slackIds ...string) Permissions
	FilterByLinks(links ...string) Permissions
	FilterByGreaterTime(time time.Time) Permissions
	FilterByLowerTime(time time.Time) Permissions
	SearchBy(search string) Permissions

	WithUsers() Permissions
	FilterByUserIds(userIds ...int64) Permissions

	Count() Permissions
	CountWithUsers() Permissions
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) Permissions
}

type Permission struct {
	RequestId   string    `json:"request_id" db:"request_id" structs:"request_id"`
	WorkSpace   string    `json:"workspace" db:"workspace" structs:"workspace"`
	SlackId     string    `json:"slack_id" db:"slack_id" structs:"slack_id"`
	Username    string    `json:"username" db:"username" structs:"username"`
	AccessLevel string    `json:"access_level" db:"access_level" structs:"access_level"`
	Link        string    `json:"link" db:"link" structs:"link"` //mean conversation
	SubmoduleId string    `json:"submodule_id" db:"submodule_id" structs:"submodule_id"`
	Bill        bool      `json:"bill" db:"bill" structs:"bill"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" structs:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" structs:"-"`
	*User       `structs:",omitempty"`
}
