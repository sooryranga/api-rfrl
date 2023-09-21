package document

import (
	"crypto/rsa"

	"github.com/jmoiron/sqlx"
)

// Handler handles all request within auth
type Handler struct {
	db  *sqlx.DB
	key *rsa.PublicKey
}

// NewHandler creates a handler
func NewHandler(db *sqlx.DB, key *rsa.PublicKey) *Handler {
	return &Handler{
		db:  db,
		key: key,
	}
}
