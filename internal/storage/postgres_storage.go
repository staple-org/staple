package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"

	"github.com/staple-org/staple/internal/models"
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
func (p PostgresStorer) Create(staple models.Staple, userID string) error {
	conn, err := connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "insert into staples(name, id, content, archived, created_timestamp, user_id) values($1, $2, $3, $4, $5, $6)",
		staple.Name,
		staple.ID,
		staple.Content,
		staple.Archived,
		staple.CreatedTimestamp,
		userID)
	return err
}

// Delete removes a staple.
func (p PostgresStorer) Delete(userID string, stapleID string) error {
	panic("implement me")
}

// Get retrieves a staple.
func (p PostgresStorer) Get(userID string, stapleID string) (models.Staple, error) {
	panic("implement me")
}

// List gets all the not archived staples for a user. List will not retrieve the content
// since that can possibly be a large text. We only ever retrieve it when that
// specific staple is Get.
func (p PostgresStorer) List(userID string) ([]models.Staple, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "select name, id, archived, created_timestamp from staples where user_id=$1 and archived = false", userID)
	if err != nil {
		return nil, err
	}
	ret := make([]models.Staple, 0)
	for rows.Next() {
		staple := models.Staple{}
		err = rows.Scan(&staple.Name, &staple.ID, &staple.Archived, &staple.CreatedTimestamp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, staple)
	}
	return ret, nil
}
