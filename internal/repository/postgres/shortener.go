package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/AntonChubarov/goit-cloud-fp/internal/model"
	"github.com/jmoiron/sqlx"
)

type LinkRepo struct {
	db *sqlx.DB
}

func NewLinkRepo(db *sqlx.DB) *LinkRepo {
	slog.Debug("Postgres link repository initialized")
	return &LinkRepo{db: db}
}

func (r *LinkRepo) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	var l model.Link
	slog.Debug("fetching link by short code", "code", code)

	err := r.db.GetContext(ctx, &l, `SELECT * FROM links WHERE short_code=$1`, code)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Warn("short code not found in database", "code", code)
		return nil, nil
	}
	if err != nil {
		slog.Error("unable to query link by code", "code", code, "error", err)
		return nil, err
	}

	slog.Debug("link fetched successfully", "id", l.ID, "url", l.OriginalURL)
	return &l, nil
}

func (r *LinkRepo) Create(ctx context.Context, l *model.Link) error {
	slog.Debug("inserting new link", "id", l.ID, "short", l.ShortCode)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO links(id, short_code, original_url) VALUES ($1,$2,$3)`,
		l.ID, l.ShortCode, l.OriginalURL)

	if err != nil {
		slog.Error("unable to insert link into database", "id", l.ID, "short", l.ShortCode, "error", err)
		return err
	}

	slog.Debug("link inserted successfully", "id", l.ID)
	return nil
}

func (r *LinkRepo) IncrementClicks(ctx context.Context, id string) error {
	slog.Debug("incrementing click count", "id", id)

	_, err := r.db.ExecContext(ctx,
		`UPDATE links SET clicks = clicks + 1 WHERE id=$1`, id)

	if err != nil {
		slog.Warn("unable to increment click count", "id", id, "error", err)
		return err
	}

	slog.Debug("click count incremented", "id", id)
	return nil
}
