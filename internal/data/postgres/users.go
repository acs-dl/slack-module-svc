package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	usersTableName       = "users"
	usersIdColumn        = usersTableName + ".id"
	usersUsernameColumn  = usersTableName + ".username"
	usersRealnameColumn  = usersTableName + ".real_name"
	usersSlackIdColumn   = usersTableName + ".slack_id"
	usersCreatedAtColumn = usersTableName + ".created_at"
	usersUpdatedAtColumn = usersTableName + ".updated_at"
)

type UsersQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
	updateBuilder sq.UpdateBuilder
}

func (q UsersQ) applyFilter(filter sq.Sqlizer) UsersQ {
	return UsersQ{
		db:            q.db,
		selectBuilder: q.selectBuilder.Where(filter),
		deleteBuilder: q.deleteBuilder.Where(filter),
		updateBuilder: q.updateBuilder.Where(filter),
	}
}

var (
	usersColumns = []string{
		usersIdColumn,
		usersUsernameColumn,
		usersRealnameColumn,
		usersSlackIdColumn,
		usersCreatedAtColumn,
	}
	selectedUsersTable = sq.Select("*").From(usersTableName)
)

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:            db.Clone(),
		selectBuilder: selectedUsersTable,
		deleteBuilder: sq.Delete(usersTableName),
		updateBuilder: sq.Update(usersTableName),
	}
}

func (q UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q UsersQ) Upsert(user data.User) (*int64, error) {
	if user.Username != nil && *user.Username == "" {
		user.Username = nil
	}

	clauses := structs.Map(user)
	updateQuery := sq.Update(" ").
		Set("username", user.Username).
		Set("updated_at", time.Now())

	if user.Id != nil {
		updateQuery = updateQuery.Set("id", *user.Id)
	}

	updateStmt, args := updateQuery.MustSql()
	query := sq.Insert(usersTableName).
		SetMap(clauses).
		Suffix("ON CONFLICT (slack_id) DO "+updateStmt, args...).
		Suffix("RETURNING *")

	var response []data.User
	err := q.db.Select(&response, query)
	if len(response) == 0 {
		return nil, errors.Wrap(err, "failed to get upsert response")
	}

	return response[0].Id, nil
}

func (q UsersQ) Delete() error {
	var deleted []data.User

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return errors.Wrap(err, "failed to delete users")
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, errors.Wrap(err, "failed to get user")
}

func (q UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.selectBuilder)

	return result, errors.Wrap(err, "failed to select users")
}

func (q UsersQ) FilterById(id *int64) data.Users {
	equalId := sq.Eq{usersIdColumn: id}

	return q.applyFilter(equalId)
}

func (q UsersQ) FilterBySlackIds(slackIds ...string) data.Users {

	equalSlackIds := sq.Eq{usersSlackIdColumn: slackIds}

	return q.applyFilter(equalSlackIds)
}

func (q UsersQ) FilterByUsername(username string) data.Users {

	equalUsername := sq.Eq{usersUsernameColumn: username}

	return q.applyFilter(equalUsername)
}

func (q UsersQ) Page(pageParams pgdb.OffsetPageParams) data.Users {
	q.selectBuilder = pageParams.ApplyTo(q.selectBuilder, "username")

	return q
}

func (q UsersQ) SearchBy(search string) data.Users {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")

	q.selectBuilder = q.selectBuilder.Where(sq.ILike{usersUsernameColumn: search})

	return q
}

func (q UsersQ) Count() data.Users {
	q.selectBuilder = sq.Select("COUNT (*)").From(usersTableName)

	return q
}

func (q UsersQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.selectBuilder)

	return count, errors.Wrap(err, "failed to get total count of users")
}

func (q UsersQ) FilterByLowerTime(time time.Time) data.Users {
	lowerTime := sq.Lt{usersUpdatedAtColumn: time}

	return q.applyFilter(lowerTime)
}
