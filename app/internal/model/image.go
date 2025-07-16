package model

import (
	"database/sql"
	"time"
)

type ImageMeta struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	URL       string    `json:"url"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

type SQLImageMeta struct {
	ID        sql.NullString `db:"id"`
	FileName  sql.NullString `db:"file_name"`
	URL       sql.NullString `db:"url"`
	FileSize  sql.NullInt64  `db:"file_size"`
	CreatedAt sql.NullTime   `db:"created_at"`
}

type ImageSearchParams struct {
	ID       string
	FileName string
}
