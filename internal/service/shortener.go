package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/AntonChubarov/goit-cloud-fp/internal/model"
	"strings"

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

func NewShortener(r LinkRepository) *Shortener { return &Shortener{repo: r} }

func (s *Shortener) Create(ctx context.Context, original string) (*model.Link, error) {
	code, err := genCode(6)
	if err != nil {
		return nil, err
	}
	link := &model.Link{
		ID:          uuid.NewString(),
		ShortCode:   code,
		OriginalURL: original,
	}
	return link, s.repo.Create(ctx, link)
}

func (s *Shortener) Resolve(ctx context.Context, code string) (string, error) {
	l, err := s.repo.GetByCode(ctx, code)
	if err != nil || l == nil {
		return "", err
	}
	_ = s.repo.IncrementClicks(ctx, l.ID) // fire-and-forget
	return l.OriginalURL, nil
}

func genCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")[:n], nil
}
