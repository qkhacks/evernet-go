package actor

import (
	"context"
	"database/sql"
)

type DataStore struct {
	db *sql.DB
}

func NewDataStore(db *sql.DB) *DataStore {
	return &DataStore{db: db}
}

func (d *DataStore) Insert(ctx context.Context, actor *Actor) (*Actor, error) {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO actors (identifier, display_name, type, password, node_identifier, creator, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		actor.Identifier,
		actor.DisplayName,
		actor.Type,
		actor.Password,
		actor.NodeIdentifier,
		actor.Creator,
		actor.CreatedAt,
		actor.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return actor, nil
}

func (d *DataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM actors WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
