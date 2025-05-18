package model

import "time"

type Link struct {
	ID          string    `db:"id"`
	ShortCode   string    `db:"short_code"`
	OriginalURL string    `db:"original_url"`
	CreatedAt   time.Time `db:"created_at"`
	Clicks      int64     `db:"clicks"`
}
