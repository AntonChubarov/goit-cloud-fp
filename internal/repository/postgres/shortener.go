package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AntonChubarov/goit-cloud-fp/internal/model"

	"github.com/jmoiron/sqlx"
)

type LinkRepo struct {
	db *sqlx.DB
}

func NewLinkRepo(db *sqlx.DB) *LinkRepo { return &LinkRepo{db: db} }

func (r *LinkRepo) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	var l model.Link
	err := r.db.GetContext(ctx, &l, `SELECT * FROM links WHERE short_code=$1`, code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &l, err
}

func (r *LinkRepo) Create(ctx context.Context, l *model.Link) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO links(id, short_code, original_url) VALUES ($1,$2,$3)`,
		l.ID, l.ShortCode, l.OriginalURL)
	return err
}

func (r *LinkRepo) IncrementClicks(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE links SET clicks = clicks + 1 WHERE id=$1`, id)
	return err
}
