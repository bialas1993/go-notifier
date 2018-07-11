package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/bialas1993/go-notifier/schema"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) Insert(ctx context.Context, notify schema.Notify) error {
	_, err := r.db.Exec(
		"INSERT INTO notifications(id, title, body, service, created_at) VALUES($1, $2, $3, $4, $5)",
		notify.ID, notify.Title, notify.Body, notify.Service, notify.CreatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, skip uint64, take uint64) ([]schema.Notify, error) {
	rows, err := r.db.Query("SELECT * FROM notifications ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []schema.Notify{}
	for rows.Next() {
		notify := schema.Notify{}
		if err = rows.Scan(&notify.ID, &notify.Title, &notify.Body, &notify.Service, &notify.CreatedAt); err == nil {
			notifications = append(notifications, notify)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
