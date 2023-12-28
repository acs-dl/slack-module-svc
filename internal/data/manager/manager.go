package manager

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Manager struct {
	Db *pgdb.DB

	Responses     data.Responses
	Permissions   data.Permissions
	Users         data.Users
	Links         data.Links
	Conversations data.Conversations
}

func NewManager(db *pgdb.DB) *Manager {
	return &Manager{
		Db:            db,
		Responses:     postgres.NewResponsesQ(db),
		Permissions:   postgres.NewPermissionsQ(db),
		Users:         postgres.NewUsersQ(db),
		Links:         postgres.NewLinksQ(db),
		Conversations: postgres.NewConversationsQ(db),
	}
}

func (m *Manager) Transaction(fn func() error) error {
	return m.Db.Transaction(fn)
}
