package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/pkg/config"
)

// PostgresStapleStorer is a storer which uses Postgres as a storage backend.
type PostgresStapleStorer struct{}

// NewPostgresStapleStorer creates a new Postgres storage medium.
func NewPostgresStapleStorer() PostgresStapleStorer {
	return PostgresStapleStorer{}
}

func (p PostgresStapleStorer) connect() (*pgx.Conn, error) {
	url := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", config.Opts.Database.Hostname, config.Opts.Database.Database, config.Opts.Database.Username, config.Opts.Database.Password)
	conn, err := pgx.Connect(context.Background(), url)
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
	_, err = conn.Exec(ctx, "insert into staples(name, content, archived, created_at, user_email) values($1, $2, $3, $4, $5)",
		staple.Name,
		staple.Content,
		staple.Archived,
		staple.CreatedAt,
		email)
	return err
}

// Delete removes a staple.
func (p PostgresStapleStorer) Delete(email string, stapleID int) error {
	conn, err := p.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "delete from staples where id = $1 and user_email = $2", stapleID, email)
	return err
}

// Get retrieves a staple.
func (p PostgresStapleStorer) Get(email string, stapleID int) (*models.Staple, error) {
	conn, err := p.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	var (
		id            int
		name, content string
		archived      bool
		createdAt     time.Time
	)
	err = conn.QueryRow(ctx, "select name, id, content, archived, created_at from staples where user_email = $1 and id = $2", email, stapleID).Scan(
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
		Name:      name,
		ID:        id,
		Content:   content,
		Archived:  archived,
		CreatedAt: createdAt,
	}, nil
}

// Oldest will get the oldest staple that is not archived.
func (p PostgresStapleStorer) Oldest(email string) (*models.Staple, error) {
	conn, err := p.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	var (
		id            int
		name, content string
		archived      bool
		createdAt     time.Time
	)
	err = conn.QueryRow(ctx, "select name, id, content, archived, created_at from staples s1 where created_at = (select MIN(created_at) from staples s2 where s2.id = s1.id and s2.user_email = $1 and s2.archived = false)", email).Scan(
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
		Name:      name,
		ID:        id,
		Content:   content,
		Archived:  archived,
		CreatedAt: createdAt,
	}, nil
}

// Archive archives a staple.
func (p PostgresStapleStorer) Archive(email string, stapleID int) error {
	conn, err := p.connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "update staples set archived = true where user_email = $1 and id = $2", email, stapleID)
	return err
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
	rows, err := conn.Query(ctx, "select name, id, archived, created_at from staples where user_email=$1 and archived = false", email)
	if err != nil {
		return nil, err
	}
	ret := make([]models.Staple, 0)
	for rows.Next() {
		staple := models.Staple{}
		err = rows.Scan(&staple.Name, &staple.ID, &staple.Archived, &staple.CreatedAt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, staple)
	}
	return ret, nil
}

// ShowArchive will return the users archived staples ordered by id.
func (p PostgresStapleStorer) ShowArchive(email string) ([]models.Staple, error) {
	conn, err := p.connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "select name, id, archived, created_at from staples where user_email=$1 and archived = true order by id", email)
	if err != nil {
		return nil, err
	}
	ret := make([]models.Staple, 0)
	for rows.Next() {
		staple := models.Staple{}
		err = rows.Scan(&staple.Name, &staple.ID, &staple.Archived, &staple.CreatedAt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, staple)
	}
	return ret, nil
}
