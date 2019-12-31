package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

// PostgresStorer is a storer which uses Postgres as a storage backend.
type PostgresStorer struct{}

// NewPostgresStorer creates a new Postgres storage medium.
func NewPostgresStorer() PostgresStorer {
	return PostgresStorer{}
}

func connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("STAPLE_DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Create will create a staple in the underlying postgres storage medium.
func (p PostgresStorer) Create(userID string) error {
	conn, err := connect()
	if err != nil {
		return err
	}
}

// Delete removes a staple.
func (p PostgresStorer) Delete(userID string, stapleID string) error {
	panic("implement me")
}

// Get retrieves a staple.
func (p PostgresStorer) Get(userID string, stapleID string) ([]byte, error) {
	panic("implement me")
}

// List gets all the not archived staples for a user.
func (p PostgresStorer) List(userID string) ([]byte, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "select * from staples where user_id=$1 and archived = false", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan()
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}
