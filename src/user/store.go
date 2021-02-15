package user

import "github.com/jmoiron/sqlx"

// Store stores db instance
type Store struct {
	db *sqlx.DB
}

// NewStore creates auth store for querying
func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s Store) getUserFromId(id string) User {

}
