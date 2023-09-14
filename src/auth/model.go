package auth

import (
	"database/sql"
	"time"
)

// Auth model
type Auth struct {
	ID           int            `db:"id"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	Token        sql.NullString `db:"token"`
	AuthType     string         `db:"auth_type"`
	Email        sql.NullString `db:"email"`
	PasswordHash []byte         `db:"password_hash"`
}
