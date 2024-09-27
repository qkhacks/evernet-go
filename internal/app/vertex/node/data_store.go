package node

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

func (d *DataStore) Insert(ctx context.Context, node *Node) (*Node, error) {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO nodes (identifier, display_name, signing_private_key, signing_public_key, creator, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		node.Identifier,
		node.DisplayName,
		node.SigningPrivateKey,
		node.SigningPublicKey,
		node.Creator,
		node.CreatedAt,
		node.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return node, nil
}

func (d *DataStore) ExistsByIdentifier(ctx context.Context, identifier string) (bool, error) {
	var count int64
	err := d.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM nodes WHERE identifier = ?", identifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
