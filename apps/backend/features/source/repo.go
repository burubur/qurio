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

func (r *PostgresRepo) Save(ctx context.Context, src *Source) error {
	query := `INSERT INTO sources (url) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, src.URL).Scan(&src.ID)
}
