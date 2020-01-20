package storage

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"

	"github.com/staple-org/staple/internal/models"
)

// PostgresUserStorer is a storer which uses Postgres as a storage backend.
type PostgresUserStorer struct{}

// NewPostgresUserStorer creates a new Postgres storage medium.
func NewPostgresUserStorer() PostgresUserStorer {
	return PostgresUserStorer{}
}

// Create saves a user in the db.
func (s PostgresUserStorer) Create(email string, password []byte) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "insert into users(email, password) values($1, $2)",
		email,
		password)
	return err
}

// Delete deletes a user from the db.
func (s PostgresUserStorer) Delete(email string) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "delete from users where email = $1",
		email)
	return err
}

// Get retrieves a user.
func (s PostgresUserStorer) Get(email string) (*models.User, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	var (
		storedEmail string
		password    []byte
	)
	err = conn.QueryRow(ctx, "select email, password from users where email = $1", email).Scan(&storedEmail, &password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &models.User{Email: storedEmail, Password: string(password)}, nil
}

// Update updates a user with a given email address.
func (s PostgresUserStorer) Update(email string, newUser models.User) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "update users set email=$1, password=$2 where email=$3",
		newUser.Email,
		newUser.Password,
		email)
	return err
}

func (s PostgresUserStorer) connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("STAPLE_DATABASE_URL"))
	if err != nil {
		log.Println("Failed to connect to database: ", err)
		return nil, err
	}
	return conn, nil
}
