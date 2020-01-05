package storage

import (
	"context"
	"os"

	"github.com/staple-org/staple/internal/models"

	"github.com/jackc/pgx/v4"
)

// PostgresUserStorer is a storer which uses Postgres as a storage backend.
type PostgresUserStorer struct{}

// NewPostgresUserStorer creates a new Postgres storage medium.
func NewPostgresUserStorer() PostgresUserStorer {
	return PostgresUserStorer{}
}

// Create saves a user in the db.
func (s PostgresUserStorer) Create(email string, password []byte) error {
	panic("implement me")
}

// Delete deletes a user from the db.
func (s PostgresUserStorer) Delete(email string) error {
	panic("implement me")
}

// Get retrieves a user.
func (s PostgresUserStorer) Get(email string) (*models.User, error) {
	panic("implement me")
}

// Update updates a user with a given email address.
func (s PostgresUserStorer) Update(email string, newUser models.User) error {
	panic("implement me")
}

func (s PostgresUserStorer) connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("STAPLE_DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
