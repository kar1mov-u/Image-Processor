package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID         uuid.UUID      `db:"id"`
	UserID     uuid.UUID      `db:"user_id"`
	Name       string         `db:"name"`
	Location   sql.NullString `db:"location"`
	UpdatedAt  time.Time      `db:"updated_at"`
	UploadedAt time.Time      `db:"uploaded_at"`
}
