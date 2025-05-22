package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"strings"
	"time"

	"github.com/AntonChubarov/goit-cloud-fp/internal/model"
	"github.com/google/uuid"
)

type LinkRepository interface {
	GetByCode(ctx context.Context, code string) (*model.Link, error)
	Create(ctx context.Context, l *model.Link) error
	IncrementClicks(ctx context.Context, id string) error
}

type Shortener struct {
	repo LinkRepository
}

func NewShortener(r LinkRepository) *Shortener {
	slog.Debug("shortener service initialized")
	return &Shortener{repo: r}
}

func (s *Shortener) Create(ctx context.Context, original string) (*model.Link, error) {
	code, err := genCode(6)
	if err != nil {
		slog.Error("unable to generate short code", "error", err)
		return nil, err
	}
	slog.Debug("short code generated", "code", code)

	link := &model.Link{
		ID:          uuid.NewString(),
		ShortCode:   code,
		OriginalURL: original,
	}

	err = s.repo.Create(ctx, link)
	if err != nil {
		slog.Error("unable to create link in repository", "short", code, "url", original, "error", err)
		return nil, err
	}
	slog.Debug("link created in repository", "id", link.ID, "short", link.ShortCode)

	return link, nil
}

func (s *Shortener) Resolve(ctx context.Context, code string) (string, error) {
	slog.Debug("resolving short code", "code", code)

	l, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		slog.Error("unable to fetch link from repository", "code", code, "error", err)
		return "", err
	}
	if l == nil {
		slog.Warn("short code not found", "code", code)
		return "", nil
	}

	// Fire-and-forget
	go func(id string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if incrementErr := s.repo.IncrementClicks(bgCtx, id); incrementErr != nil {
			slog.Warn("unable to increment click count", "id", id, "error", incrementErr)
		} else {
			slog.Debug("click count incremented", "id", id)
		}
	}(l.ID)

	return l.OriginalURL, nil
}

func genCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")[:n]
	return code, nil
}
