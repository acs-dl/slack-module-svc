package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	linksTableName  = "links"
	linksLinkColumn = linksTableName + ".link"
)

type LinksQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

func (r LinksQ) applyFilter(filter sq.Sqlizer) LinksQ {
	return LinksQ{
		db:            r.db,
		selectBuilder: r.selectBuilder.Where(filter),
		deleteBuilder: r.deleteBuilder.Where(filter),
	}
}

func NewLinksQ(db *pgdb.DB) data.Links {
	return &LinksQ{
		db:            db,
		selectBuilder: sq.Select(linksTableName + ".*").From(linksTableName),
		deleteBuilder: sq.Delete(linksTableName),
	}
}

func (r LinksQ) New() data.Links {
	return NewLinksQ(r.db)
}

func (r LinksQ) Get() (*data.Link, error) {
	var result data.Link
	err := r.db.Get(&result, r.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, errors.Wrap(err, "failed to get links")
}

func (r LinksQ) Select() ([]data.Link, error) {
	var result []data.Link

	err := r.db.Select(&result, r.selectBuilder)

	return result, errors.Wrap(err, "failed to select links")
}

func (r LinksQ) Insert(link data.Link) error {
	insertStmt := sq.Insert(linksTableName).SetMap(structs.Map(link)).Suffix("ON CONFLICT (link) DO NOTHING")
	err := r.db.Exec(insertStmt)
	
	return errors.Wrap(err, "failed to insert link")
}

func (r LinksQ) Delete() error {
	var deleted []data.Link

	err := r.db.Select(&deleted, r.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return errors.Wrap(err, "failed to delete links")
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r LinksQ) FilterByLinks(links ...string) data.Links {
	equalLinks := sq.Eq{linksLinkColumn: links}

	return r.applyFilter(equalLinks)
}
