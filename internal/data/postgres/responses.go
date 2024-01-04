package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	responsesTableName = "responses"
	responsesIdColumn  = responsesTableName + ".id"
)

type ResponsesQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

func (q ResponsesQ) applyFilter(filter sq.Sqlizer) ResponsesQ {
	return ResponsesQ{
		db:            q.db,
		selectBuilder: q.selectBuilder.Where(filter),
		deleteBuilder: q.deleteBuilder.Where(filter),
	}
}

var selectedResponsesTable = sq.Select("*").From(responsesTableName)

func NewResponsesQ(db *pgdb.DB) data.Responses {
	return &ResponsesQ{
		db:            db.Clone(),
		selectBuilder: selectedResponsesTable,
		deleteBuilder: sq.Delete(responsesTableName),
	}
}

func (q ResponsesQ) New() data.Responses {
	return NewResponsesQ(q.db)
}

func (q ResponsesQ) Insert(response data.Response) error {
	clauses := structs.Map(response)

	query := sq.Insert(responsesTableName).SetMap(clauses)

	return errors.Wrap(q.db.Exec(query), "failed to insert response")
}

func (q ResponsesQ) Select() ([]data.Response, error) {
	var result []data.Response

	err := q.db.Select(&result, q.selectBuilder)

	return result, errors.Wrap(err, "failed to select responses")
}

func (q ResponsesQ) Get() (*data.Response, error) {
	var result data.Response

	err := q.db.Get(&result, q.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, errors.Wrap(err, "failed to get response")
}

func (q ResponsesQ) Delete() error {
	var deleted []data.Response

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return errors.Wrap(err, "failed to delete responses")
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q ResponsesQ) FilterByIds(ids ...string) data.Responses {
	equalIds := sq.Eq{responsesIdColumn: ids}

	return q.applyFilter(equalIds)
}
