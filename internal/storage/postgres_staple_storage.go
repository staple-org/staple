package storage

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/staple-org/staple/internal/models"
)

// PostgresStapleStorer is a storer which uses Postgres as a storage backend.
type PostgresStapleStorer struct{}

// NewPostgresStapleStorer creates a new Postgres storage medium.
func NewPostgresStapleStorer() PostgresStapleStorer {
	return PostgresStapleStorer{}
}

func (p PostgresStapleStorer) connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("STAPLE_DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Create will create a staple in the underlying postgres storage medium.
func (p PostgresStapleStorer) Create(staple models.Staple, email string) error {
	conn, err := p.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "insert into staples(name, id, content, archived, created_timestamp, username) values($1, $2, $3, $4, $5, $6)",
		staple.Name,
		staple.ID,
		staple.Content,
		staple.Archived,
		staple.CreatedTimestamp,
		email)
	return err
}

// Delete removes a staple.
func (p PostgresStapleStorer) Delete(email string, stapleID string) error {
	panic("implement me")
}

// Get retrieves a staple.
func (p PostgresStapleStorer) Get(email string, stapleID string) (*models.Staple, error) {
	conn, err := p.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	var (
		name, id, content string
		archived          bool
		createdAt         time.Time
	)
	err = conn.QueryRow(ctx, "select name, id, content, archived, created_timestamp from staples where user_email = $1 and id = $2", email, stapleID).Scan(
		&name,
		&id,
		&content,
		&archived,
		&createdAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &models.Staple{
		Name:             name,
		ID:               id,
		Content:          content,
		Archived:         archived,
		CreatedTimestamp: createdAt,
	}, nil
}

// Archive archives a staple.
func (p PostgresStapleStorer) Archive(email string, stapleID string) error {
	panic("implement me")
}

// List gets all the not archived staples for a user. List will not retrieve the content
// since that can possibly be a large text. We only ever retrieve it when that
// specific staple is Get.
func (p PostgresStapleStorer) List(email string) ([]models.Staple, error) {
	conn, err := p.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "select name, id, archived, created_timestamp from staples where user_email=$1 and archived = false", email)
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
