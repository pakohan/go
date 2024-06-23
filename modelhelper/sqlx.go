package modelhelper

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func (db *DB) NamedGetContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	return db.Transaction(ctx, nil, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			return err
		}

		return stmt.GetContext(ctx, dest, arg)
	})
}

func (db *DB) Transaction(ctx context.Context, opts *sql.TxOptions, f func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	err = f(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}
