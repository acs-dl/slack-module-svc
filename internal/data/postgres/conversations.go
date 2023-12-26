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
	conversationsTableName   = "conversations"
	conversationsTitleColumn = conversationsTableName + ".title"
	conversationsIdColumn    = conversationsTableName + ".id"
)

type ConversationsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

func (r ConversationsQ) applyFilter(filter sq.Sqlizer) ConversationsQ {
	return ConversationsQ{
		db:            r.db,
		selectBuilder: r.selectBuilder.Where(filter),
		deleteBuilder: r.deleteBuilder.Where(filter),
	}
}

func NewConversationsQ(db *pgdb.DB) data.Conversations {
	return &ConversationsQ{
		db:            db,
		selectBuilder: sq.Select(conversationsTableName + ".*").From(conversationsTableName),
		deleteBuilder: sq.Delete(conversationsTableName),
	}
}

func (r ConversationsQ) New() data.Conversations {
	return NewConversationsQ(r.db)
}

func (r ConversationsQ) Get() (*data.Conversation, error) {
	var result data.Conversation
	err := r.db.Get(&result, r.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, errors.Wrap(err, "failed to get conversation")
}

func (r ConversationsQ) Select() ([]data.Conversation, error) {
	var result []data.Conversation

	err := r.db.Select(&result, r.selectBuilder)

	return result, errors.Wrap(err, "failed to select conversations")
}

func (r ConversationsQ) Upsert(conversation data.Conversation) error {
	updateStmt, args := sq.Update(" ").
		Set("title", conversation.Title).
		Set("members_amount", conversation.MembersAmount).
		MustSql()

	query := sq.Insert(conversationsTableName).SetMap(structs.Map(conversation)).
		Suffix("ON CONFLICT (id) DO "+updateStmt, args...)

	err := r.db.Exec(query)
	return errors.Wrap(err, "failed to insert conversation")
}

func (r ConversationsQ) Delete() error {
	var deleted []data.Conversation

	err := r.db.Select(&deleted, r.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return errors.Wrap(err, "failed to delete conversations")
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r ConversationsQ) FilterByTitles(titles ...string) data.Conversations {
	equalTitles := sq.Eq{conversationsTitleColumn: titles}

	return r.applyFilter(equalTitles)
}

func (r ConversationsQ) FilterByIds(ids ...string) data.Conversations {
	equalIds := sq.Eq{conversationsIdColumn: ids}

	return r.applyFilter(equalIds)
}
