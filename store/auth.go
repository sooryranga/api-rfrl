package store

import (
	"errors"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"gopkg.in/guregu/null.v4"
)

// AuthStore holds all store related functions for client
type AuthStore struct{}

// NewAuthStore creates new authstore
func NewAuthStore() *AuthStore {
	return &AuthStore{}
}

const checkEmailExists string = `
SELECT EXISTS (
	SELECT 1 FROM auth 
	WHERE email=$1 AND auth.auth_type = 'email' AND client_id != $2
)
`

func (au *AuthStore) CheckEmailAuthExists(db tutorme.DB, clientID string, email string) (bool, error) {
	var exists null.Bool
	err := db.QueryRowx(checkEmailExists, email, clientID).Scan(&exists)

	if !exists.Valid {
		return false, errors.New("Exists returned a non bool value")
	}
	return exists.Bool, err
}

const updateEmail string = `
UPDATE auth 
SET auth.email = $1
WHERE client_id = $2
AND auth_type = 'email'
	`

func (au *AuthStore) UpdateAuthEmail(db tutorme.DB, clientID string, email string) error {
	_, err := db.Queryx(updateEmail, email, clientID)

	return err
}

const getByToken string = `
SELECT client.* FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.token =$1 AND auth.auth_type =$2 
LIMIT 1
	`

// GetByToken queries the database for token auth from providers
func (au *AuthStore) GetByToken(db tutorme.DB, token string, authType string) (*tutorme.Client, error) {
	var c tutorme.Client
	err := db.QueryRowx(getByToken, token, authType).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

const getByEmail string = `
SELECT password_hash, client.* FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.email =$1 AND auth.auth_type =$2 
LIMIT 1
	`

// GetByEmail queries the database for email auth
func (au *AuthStore) GetByEmail(db tutorme.DB, email string) (*tutorme.Client, []byte, error) {

	type getByEmailStruct struct {
		tutorme.Client
		PasswordHash []byte `db:"password_hash"`
	}
	var result getByEmailStruct

	row := db.QueryRowx(getByEmail, email, tutorme.EMAIL)

	err := row.StructScan(&result)

	return &result.Client, result.PasswordHash, err
}

const insertEmailAuth string = `
INSERT INTO auth (email, password_hash, auth_type, client_id)
VALUES ($1, $2, $3, $4)
RETURNING id
	`

// CreateWithEmail creates auth row with email in db
func (au *AuthStore) CreateWithEmail(db tutorme.DB, auth *tutorme.Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertEmailAuth, auth.Email, auth.PasswordHash, tutorme.EMAIL, clientID)
	var id int = -1

	err := row.Scan(&id)

	return id, err
}

const insertToken string = `
INSERT INTO auth (token, auth_type, client_id) 
VALUES ($1, $2, $3) 
RETURNING id
	`

// CreateWithToken creates auth row with token in db
func (au *AuthStore) CreateWithToken(db tutorme.DB, auth *tutorme.Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertToken, auth.Token, auth.AuthType, clientID)
	var id int = -1
	err := row.Scan(&id)

	return id, err
}
