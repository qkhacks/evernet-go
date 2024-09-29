package messaging

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
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

func (d *InboxDataStore) FindByActorAddressAndNodeIdentifier(ctx context.Context, actorAddress string, nodeIdentifier string, page int64, size int64) ([]*Inbox, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT identifier, display_name, node_identifier, actor_address, created_at, updated_at FROM inboxes WHERE actor_address = ? AND node_identifier = ? LIMIT ? OFFSET ?",
		actorAddress, nodeIdentifier, size, page*size)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zap.L().Error("error closing rows", zap.Error(err))
		}
	}(rows)

	var inboxes []*Inbox

	for rows.Next() {
		var inbox Inbox
		err = rows.Scan(&inbox.Identifier, &inbox.DisplayName, &inbox.NodeIdentifier, &inbox.ActorAddress, &inbox.CreatedAt, &inbox.UpdatedAt)

		if err != nil {
			return nil, err
		}

		inboxes = append(inboxes, &inbox)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inboxes, nil
}

func (d *InboxDataStore) FindByIdentifierAndActorAddressAndNodeIdentifier(ctx context.Context, identifier string, actorAddress string, nodeIdentifier string) (*Inbox, error) {
	var inbox Inbox
	err := d.db.QueryRowContext(ctx,
		"SELECT identifier, display_name, node_identifier, actor_address, created_at, updated_at FROM inboxes WHERE identifier = ? AND actor_address = ? AND node_identifier = ?",
		identifier, actorAddress, nodeIdentifier).Scan(
		&inbox.Identifier,
		&inbox.DisplayName,
		&inbox.NodeIdentifier,
		&inbox.ActorAddress,
		&inbox.CreatedAt,
		&inbox.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &inbox, nil
}

func (d *InboxDataStore) UpdateDisplayNameByIdentifierAndActorAddressAndNodeIdentifier(ctx context.Context, displayName string, identifier string, actorAddress string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"UPDATE inboxes SET display_name = ? WHERE identifier = ? AND actor_address = ? AND node_identifier = ?",
		displayName, identifier, actorAddress, nodeIdentifier)

	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *InboxDataStore) DeleteByIdentifierAndActorAddressAndNodeIdentifier(ctx context.Context, identifier string, actorAddress string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"DELETE FROM inboxes WHERE identifier = ? AND actor_address = ? AND node_identifier = ?",
		identifier, actorAddress, nodeIdentifier)

	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *InboxDataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inboxes WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
