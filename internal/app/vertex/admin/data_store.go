package admin

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

func (d *DataStore) Insert(ctx context.Context, a *Admin) (*Admin, error) {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO admins (identifier, password, creator, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		a.Identifier, a.Password, a.Creator, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (d *DataStore) FindByIdentifier(ctx context.Context, identifier string) (*Admin, error) {
	var a Admin

	err := d.db.QueryRowContext(ctx,
		"SELECT identifier, password, creator, created_at, updated_at FROM admins WHERE identifier = ?", identifier).
		Scan(&a.Identifier, &a.Password, &a.Creator, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (d *DataStore) FindAll(ctx context.Context, page int64, size int64) ([]*Admin, error) {
	var admins []*Admin
	rows, err := d.db.QueryContext(ctx,
		"SELECT identifier, password, creator, created_at, updated_at FROM admins LIMIT ? OFFSET ?",
		size, page*size)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zap.L().Error("error while closing rows", zap.Error(err))
		}
	}(rows)

	for rows.Next() {
		var a Admin
		err = rows.Scan(&a.Identifier, &a.Password, &a.Creator, &a.CreatedAt, &a.UpdatedAt)

		if err != nil {
			return nil, err
		}

		admins = append(admins, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}

func (d *DataStore) UpdatePasswordByIdentifier(ctx context.Context, password string, identifier string) error {
	result, err := d.db.ExecContext(ctx, "UPDATE admins SET password = ? WHERE identifier = ?", password, identifier)

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *DataStore) DeleteByIdentifier(ctx context.Context, identifier string) error {
	result, err := d.db.ExecContext(ctx, "DELETE FROM admins WHERE identifier = ?", identifier)

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *DataStore) Exists(ctx context.Context) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admins").Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d *DataStore) ExistsByIdentifier(ctx context.Context, identifier string) (bool, error) {
	var count int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admins WHERE identifier = ?", identifier).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
