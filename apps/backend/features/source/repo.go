package source

import (
	"context"
	"database/sql"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sources WHERE content_hash = $1)`
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresRepo) Save(ctx context.Context, src *Source) error {
	query := `INSERT INTO sources (url, content_hash) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, src.URL, src.ContentHash).Scan(&src.ID)
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE sources SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *PostgresRepo) List(ctx context.Context) ([]Source, error) {
	query := `SELECT id, url, status FROM sources ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var s Source
		if err := rows.Scan(&s.ID, &s.URL, &s.Status); err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, nil
}