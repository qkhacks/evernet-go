package messaging

import (
	"context"
	"database/sql"
)

type OutboxDataStore struct {
	db *sql.DB
}

func NewOutboxDataStore(db *sql.DB) *OutboxDataStore {
	return &OutboxDataStore{db: db}
}

func (d *OutboxDataStore) Insert(ctx context.Context, outbox *Outbox) (*Outbox, error) {
	_, err := d.db.ExecContext(ctx, "INSERT INTO outboxes (identifier, display_name, node_identifier, actor_address, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		outbox.Identifier,
		outbox.DisplayName,
		outbox.NodeIdentifier,
		outbox.ActorAddress,
		outbox.CreatedAt,
		outbox.UpdatedAt)

	return outbox, err
}

func (d *OutboxDataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outboxes WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
