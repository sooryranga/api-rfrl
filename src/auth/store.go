package auth

import (
	"github.com/jmoiron/sqlx"
)

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

const (
	getByToken string = `
SELECT * FROM auth 
WHERE auth.token =$1 AND auth.type =$2 
LIMIT 1
	`
	getByEmail string = `
SELECT * FROM auth 
WHERE auth.email =$1 AND auth.type =$2 
LIMIT 1
	`
	insertEmailAuth string = `
INSERT INTO auth (email, password, auth_type)
VALUES ($1, $2, $3)
RETURNING id
	`
	insertToken string = `
INSERT INTO auth (token, auth_type) 
VALUES ($1, $2) 
RETURNING id
	`
)

// GetByToken queries the database for token auth from providers
func (au *Store) GetByToken(token string, authType string) (*Auth, error) {
	var m Auth
	err := au.db.QueryRowx(getByToken, token, authType).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByEmail queries the database for email auth
func (au *Store) GetByEmail(email string) (*Auth, error) {
	var m Auth
	err := au.db.QueryRowx(getByEmail, email, EMAIL).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateWithEmail creates auth row with email in db
func (au *Store) CreateWithEmail(email string, password []byte) (int, error) {
	row := au.db.QueryRow(insertEmailAuth, email, password, EMAIL)
	var id int

	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

// CreateWithToken creates auth row with token in db
func (au *Store) CreateWithToken(token string, authType string) (int, error) {
	row := au.db.QueryRow(insertToken, token, authType)
	var id int

	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
