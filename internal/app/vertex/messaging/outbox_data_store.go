package messaging

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
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

func (d *OutboxDataStore) FindByActorAddressAndNodeIdentifier(ctx context.Context, actorAddress string, nodeIdentifier string, page int64, size int64) ([]*Outbox, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT * FROM outboxes WHERE actor_address = ? and node_identifier = ? LIMIT ? OFFSET ?", actorAddress, nodeIdentifier, size, page*size)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zap.L().Error("error closing rows", zap.Error(err))
		}
	}(rows)

	var outboxes []*Outbox
	for rows.Next() {
		var outbox Outbox
		err = rows.Scan(&outbox.Identifier, &outbox.DisplayName, &outbox.NodeIdentifier, &outbox.ActorAddress, &outbox.CreatedAt, &outbox.UpdatedAt)

		if err != nil {
			return nil, err
		}

		outboxes = append(outboxes, &outbox)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return outboxes, nil
}

func (d *OutboxDataStore) FindByIdentifierAndActorAddressAndNodeIdentifier(ctx context.Context, identifier string, actorAddress string, nodeIdentifier string) (*Outbox, error) {
	var outbox Outbox

	err := d.db.QueryRowContext(ctx,
		"SELECT identifier, display_name, node_identifier, actor_address, created_at, updated_at FROM outboxes WHERE identifier = ? AND actor_address = ? AND node_identifier = ?", identifier, actorAddress, nodeIdentifier).
		Scan(
			&outbox.Identifier,
			&outbox.DisplayName,
			&outbox.NodeIdentifier,
			&outbox.ActorAddress,
			&outbox.CreatedAt,
			&outbox.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &outbox, nil
}

func (d *OutboxDataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outboxes WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
