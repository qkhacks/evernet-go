package messaging

import (
	"context"
	"database/sql"
)

type InboxDataStore struct {
	db *sql.DB
}

func NewInboxDataStore(db *sql.DB) *InboxDataStore {
	return &InboxDataStore{db: db}
}

func (d *InboxDataStore) Insert(ctx context.Context, inbox *Inbox) (*Inbox, error) {
	_, err := d.db.ExecContext(ctx, "INSERT INTO inboxes (identifier, display_name, node_identifier, actor_address, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		inbox.Identifier,
		inbox.DisplayName,
		inbox.NodeIdentifier,
		inbox.ActorAddress,
		inbox.CreatedAt,
		inbox.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return inbox, nil
}

func (d *InboxDataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inboxes WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
