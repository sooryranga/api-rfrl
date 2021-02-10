package auth

import (
	"database/sql"
	"time"
)

// Auth model
type Auth struct {
	ID       int
	Created  time.Time
	Updated  time.Time
	Token    sql.NullString
	AuthType string
	Password sql.NullString
	Email    sql.NullString
	Salt     sql.NullString
}
