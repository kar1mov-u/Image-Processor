package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserName   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"`
	Created_at time.Time `json:"created_at" db: created_at`
}
