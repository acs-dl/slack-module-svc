package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/helpers"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	permissionsTableName         = "permissions"
	permissionsRequestIdColumn   = permissionsTableName + ".request_id"
	permissionsWorkSpaceColumn   = permissionsTableName + ".workspace"
	permissionsSlackIdColumn     = permissionsTableName + ".slack_id"
	permissionsLinkColumn        = permissionsTableName + ".link"
	permissonBillColumn          = permissionsTableName + ".bill"
	permissionsSubmoduleIdColumn = permissionsTableName + ".submodule_id"
	//permissionsSubmoduleAccessHashColumn = permissionsTableName + ".submodule_access_hash"
	permissionsAccessLevelColumn = permissionsTableName + ".access_level"
	permissionsCreatedAtColumn   = permissionsTableName + ".created_at"
	permissionsUpdatedAtColumn   = permissionsTableName + ".updated_at"
)

type PermissionsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
	updateBuilder sq.UpdateBuilder
}

func (q PermissionsQ) applyFilter(filter sq.Sqlizer) PermissionsQ {
	return PermissionsQ{
		db:            q.db,
		selectBuilder: q.selectBuilder.Where(filter),
		deleteBuilder: q.deleteBuilder.Where(filter),
		updateBuilder: q.updateBuilder.Where(filter),
	}
}

var permissionsColumns = []string{
	permissionsRequestIdColumn,
	permissionsWorkSpaceColumn,
	permissionsSlackIdColumn,
	permissionsLinkColumn,
	permissonBillColumn,
	permissionsAccessLevelColumn,
	permissionsCreatedAtColumn,
	permissionsUpdatedAtColumn,
	permissionsSubmoduleIdColumn,
}

func NewPermissionsQ(db *pgdb.DB) data.Permissions {
	return &PermissionsQ{
		db:            db.Clone(),
		selectBuilder: sq.Select(permissionsColumns...).From(permissionsTableName),
		deleteBuilder: sq.Delete(permissionsTableName),
		updateBuilder: sq.Update(permissionsTableName),
	}
}

func (q PermissionsQ) New() data.Permissions {
	return NewPermissionsQ(q.db)
}

func (q PermissionsQ) UpdateAccessLevel(permission data.Permission) error {
	query := q.updateBuilder.Set("access_level", permission.AccessLevel)

	return q.db.Exec(query)
}

func (q PermissionsQ) Select() ([]data.Permission, error) {
	var result []data.Permission

	err := q.db.Select(&result, q.selectBuilder)

	return result, errors.Wrap(err, "failed to select permissions")
}

func (q PermissionsQ) Upsert(permission data.Permission) error {
	updateStmt, args := sq.Update(" ").
		Set("updated_at", time.Now()).
		Set("bill", permission.Bill).
		Set("access_level", permission.AccessLevel).
		Set("request_id", permission.RequestId).
		MustSql()

	query := sq.Insert(permissionsTableName).SetMap(structs.Map(permission)).
		Suffix("ON CONFLICT (slack_id, submodule_id) DO "+updateStmt, args...)

	return errors.Wrap(q.db.Exec(query), "failed to insert permission")
}

func (q PermissionsQ) Get() (*data.Permission, error) {
	var result data.Permission

	err := q.db.Get(&result, q.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, errors.Wrap(err, "failed to get permission")
}

func (q PermissionsQ) Delete() error {
	var deleted []data.Permission

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return errors.Wrap(err, "failed to delete permissions")
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q PermissionsQ) FilterBySlackIds(slackIds ...string) data.Permissions {
	equalSlackIds := sq.Eq{permissionsSlackIdColumn: slackIds}

	return q.applyFilter(equalSlackIds)
}

func (q PermissionsQ) FilterByLinks(links ...string) data.Permissions {
	equalLinks := sq.Eq{permissionsLinkColumn: links}

	return q.applyFilter(equalLinks)
}

func (q PermissionsQ) SearchBy(search string) data.Permissions {
	search = strings.Replace(search, " ", "%", -1)
	search = fmt.Sprint("%", search, "%")
	ilikeSearch := sq.ILike{permissionsLinkColumn: search}

	return q.applyFilter(ilikeSearch)
}

func (q PermissionsQ) Count() data.Permissions {
	q.selectBuilder = sq.Select("COUNT (*)").From(permissionsTableName)

	return q
}

func (q PermissionsQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.selectBuilder)

	return count, errors.Wrap(err, "failed to get total count")
}

func (q PermissionsQ) Page(pageParams pgdb.OffsetPageParams) data.Permissions {
	q.selectBuilder = pageParams.ApplyTo(q.selectBuilder, "link")

	return q
}

func (q PermissionsQ) WithUsers() data.Permissions {
	q.selectBuilder = sq.Select().Columns(helpers.RemoveDuplicateColumn(append(permissionsColumns, usersColumns...))...).
		From(permissionsTableName).
		LeftJoin(usersTableName + " ON " + usersSlackIdColumn + " = " + permissionsSlackIdColumn).
		Where(sq.NotEq{permissionsRequestIdColumn: nil}).
		GroupBy(helpers.RemoveDuplicateColumn(append(permissionsColumns, usersColumns...))...)

	return q
}

func (q PermissionsQ) CountWithUsers() data.Permissions {
	q.selectBuilder = sq.Select("COUNT(*)").From(permissionsTableName).
		LeftJoin(usersTableName + " ON " + usersSlackIdColumn + " = " + permissionsSlackIdColumn).
		Where(sq.NotEq{permissionsRequestIdColumn: nil})

	return q
}

func (q PermissionsQ) FilterByUserIds(userIds ...int64) data.Permissions {
	equalUserIds := sq.Eq{usersIdColumn: userIds}

	if len(userIds) == 0 {
		equalUserIds = sq.Eq{usersIdColumn: nil}
	}

	return q.applyFilter(equalUserIds)
}

func (q PermissionsQ) FilterByGreaterTime(time time.Time) data.Permissions {
	greaterTime := sq.Gt{permissionsUpdatedAtColumn: time}

	return q.applyFilter(greaterTime)
}

func (q PermissionsQ) FilterByLowerTime(time time.Time) data.Permissions {
	lowerTime := sq.Lt{permissionsUpdatedAtColumn: time}

	return q.applyFilter(lowerTime)
}
