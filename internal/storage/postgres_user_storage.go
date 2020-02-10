package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/pkg/config"
)

const (
	// DefaultMaxStaples is 25.
	DefaultMaxStaples = 25
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
	_, err = conn.Exec(ctx, "insert into users(email, password, confirm_code, max_staples) values($1, $2, $3, $4)",
		email,
		password,
		"",
		DefaultMaxStaples)
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
		confirmCode string
		maxStaples  int
	)
	err = conn.QueryRow(ctx, "select email, password, confirm_code, max_staples from users where email = $1", email).Scan(&storedEmail, &password, &confirmCode, &maxStaples)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &models.User{
		Email:       storedEmail,
		Password:    string(password),
		ConfirmCode: confirmCode,
		MaxStaples:  maxStaples}, nil
}

// Update updates a user with a given email address.
func (s PostgresUserStorer) Update(email string, newUser models.User) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "update users set email=$1, password=$2, confirm_code=$3, max_staples=$4 where email=$5",
		newUser.Email,
		newUser.Password,
		newUser.ConfirmCode,
		newUser.MaxStaples,
		email)
	return err
}

func (s PostgresUserStorer) connect() (*pgx.Conn, error) {
	url := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", config.Opts.Database.Hostname, config.Opts.Database.Database, config.Opts.Database.Username, config.Opts.Database.Password)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		config.Opts.Logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}
	return conn, nil
}
