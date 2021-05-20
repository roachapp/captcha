package store

import (
	"context"
	pgx "github.com/jackc/pgx/v4/pgxpool"
)

// postgresStore is an internal store for captcha ids and their values.
type postgresStore struct {
	pgx *pgx.Pool
}

// NewPostgresStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewPostgresStore(ctx context.Context) Store {
	return &postgresStore{
		pgx: connectDB(ctx),
	}
}

func (pgs *postgresStore) Set(id string, digits []byte) {

}

func (pgs *postgresStore) Get(id string, clear bool) (digits []byte) {

	if clear {

	}
	return
}
