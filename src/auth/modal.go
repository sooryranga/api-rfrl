package auth

import (
	"database/sql"
	"time"
)

// Auth model
type Auth struct {
	ID           int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Token        sql.NullString
	AuthType     string
	Email        sql.NullString
	PasswordHash sql.RawBytes
}
