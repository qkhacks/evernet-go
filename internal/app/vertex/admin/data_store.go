package admin

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

func (d *DataStore) Insert(ctx context.Context, a *Admin) (*Admin, error) {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO admins (identifier, password, creator, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		a.Identifier, a.Password, a.Creator, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (d *DataStore) Exists(ctx context.Context) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admins").Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
