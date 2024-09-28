package node

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
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

func (d *DataStore) FindAll(ctx context.Context, page int64, size int64) ([]*Node, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT identifier, display_name, signing_private_key, signing_public_key, creator, created_at, updated_at FROM nodes LIMIT ? OFFSET ?",
		size, page*size)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zap.L().Error("failed to close rows", zap.Error(err))
		}
	}(rows)

	var nodes []*Node

	for rows.Next() {
		var node Node
		err = rows.Scan(&node.Identifier, &node.DisplayName, &node.SigningPrivateKey, &node.SigningPublicKey, &node.Creator, &node.CreatedAt, &node.UpdatedAt)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (d *DataStore) FindByIdentifier(ctx context.Context, identifier string) (*Node, error) {
	var node Node
	err := d.db.QueryRowContext(ctx,
		"SELECT identifier, display_name, signing_private_key, signing_public_key, creator, created_at, updated_at FROM nodes WHERE identifier = ?",
		identifier).
		Scan(
			&node.Identifier,
			&node.DisplayName,
			&node.SigningPrivateKey,
			&node.SigningPublicKey,
			&node.Creator,
			&node.CreatedAt,
			&node.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (d *DataStore) UpdateDisplayNameByIdentifier(ctx context.Context, displayName string, identifier string) error {
	result, err := d.db.ExecContext(ctx,
		"UPDATE nodes SET display_name = ? WHERE identifier = ?",
		displayName, identifier)

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

func (d *DataStore) DeleteByIdentifier(ctx context.Context, identifier string) error {
	result, err := d.db.ExecContext(ctx, "DELETE FROM nodes WHERE identifier = ?", identifier)

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

func (d *DataStore) ExistsByIdentifier(ctx context.Context, identifier string) (bool, error) {
	var count int64
	err := d.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM nodes WHERE identifier = ?", identifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
