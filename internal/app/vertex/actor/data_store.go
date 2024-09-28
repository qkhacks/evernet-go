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

func (d *DataStore) FindByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (*Actor, error) {
	var actor Actor

	err := d.db.QueryRowContext(ctx,
		"SELECT identifier, display_name, type, password, node_identifier, creator, created_at, updated_at FROM actors WHERE identifier = ? AND node_identifier = ?",
		identifier, nodeIdentifier).
		Scan(&actor.Identifier, &actor.DisplayName, &actor.Type, &actor.Password, &actor.NodeIdentifier, &actor.Creator, &actor.CreatedAt, &actor.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &actor, nil
}

func (d *DataStore) UpdatePasswordByIdentifierAndNodeIdentifier(ctx context.Context, password string, identifier string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"UPDATE actors SET password = ? WHERE identifier = ? AND node_identifier = ?",
		password, identifier, nodeIdentifier)

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

func (d *DataStore) UpdateDisplayNameByIdentifierAndNodeIdentifier(ctx context.Context, displayName string, identifier string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"UPDATE actors SET display_name = ? WHERE identifier = ? AND node_identifier = ?",
		displayName, identifier, nodeIdentifier)

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

func (d *DataStore) UpdateTypeByIdentifierAndNodeIdentifier(ctx context.Context, actorType string, identifier string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"UPDATE actors SET type = ? WHERE identifier = ? AND node_identifier = ?",
		actorType, identifier, nodeIdentifier)

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

func (d *DataStore) DeleteByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) error {
	result, err := d.db.ExecContext(ctx,
		"DELETE FROM actors WHERE identifier = ? AND node_identifier = ?",
		identifier, nodeIdentifier)

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

func (d *DataStore) ExistsByIdentifierAndNodeIdentifier(ctx context.Context, identifier string, nodeIdentifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM actors WHERE identifier = ? AND node_identifier = ?", identifier, nodeIdentifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
