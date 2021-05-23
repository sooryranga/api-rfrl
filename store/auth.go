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

const getAuthFromClientID string = `
SELECT * FROM auth WHERE client_id = $1
`

func (au *AuthStore) GetByClientID(db tutorme.DB, clientID string) (*tutorme.Auth, error) {
	var auth tutorme.Auth

	row := db.QueryRowx(getAuthFromClientID, clientID)

	err := row.StructScan(&auth)

	return &auth, err
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

const updateSignUpFlowQuery string = `
UPDATE auth
SET sign_up_flow = $1
WHERE client_id = $2
`

func (au *AuthStore) UpdateSignUpFlow(db tutorme.DB, clientID string, stage tutorme.SignUpFlow) error {
	_, err := db.Queryx(updateSignUpFlowQuery, stage, clientID)

	return err
}

const updateEmail string = `
UPDATE auth 
SET email = $1
WHERE client_id = $2
AND auth_type = 'email'
	`

func (au *AuthStore) UpdateAuthEmail(db tutorme.DB, clientID string, email string) error {
	_, err := db.Queryx(updateEmail, email, clientID)

	return err
}

const getByToken string = `
SELECT client.*, sign_up_flow FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.token =$1 AND auth.auth_type =$2 AND blocked = FALSE
LIMIT 1
	`

// GetByToken queries the database for token auth from providers
func (au *AuthStore) GetByToken(db tutorme.DB, token string, authType string) (*tutorme.Client, *tutorme.Auth, error) {
	type getByTokenStruct struct {
		tutorme.Client
		tutorme.Auth
	}
	var result getByTokenStruct
	err := db.QueryRowx(getByToken, token, authType).StructScan(&result)
	if err != nil {
		return nil, nil, err
	}
	return &result.Client, &result.Auth, nil
}

const getByEmail string = `
SELECT password_hash, sign_up_flow, client.* FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.email =$1 AND auth.auth_type =$2 AND blocked = FALSE
LIMIT 1
	`

// GetByEmail queries the database for email auth
func (au *AuthStore) GetByEmail(db tutorme.DB, email string) (*tutorme.Client, *tutorme.Auth, error) {

	type getByEmailStruct struct {
		tutorme.Client
		tutorme.Auth
	}
	var result getByEmailStruct

	row := db.QueryRowx(getByEmail, email, tutorme.EMAIL)

	err := row.StructScan(&result)

	return &result.Client, &result.Auth, err
}

const insertEmailAuth string = `
INSERT INTO auth (email, password_hash, auth_type, client_id)
VALUES ($1, $2, $3, $4)
RETURNING *
	`

// CreateWithEmail creates auth row with email in db
func (au *AuthStore) CreateWithEmail(db tutorme.DB, auth *tutorme.Auth, clientID string) (*tutorme.Auth, error) {
	var createdAuth tutorme.Auth
	row := db.QueryRowx(insertEmailAuth, auth.Email, auth.PasswordHash, tutorme.EMAIL, clientID)

	err := row.StructScan(&createdAuth)

	return &createdAuth, err
}

const insertToken string = `
INSERT INTO auth (token, auth_type, client_id) 
VALUES ($1, $2, $3) 
RETURNING *
	`

// CreateWithToken creates auth row with token in db
func (au *AuthStore) CreateWithToken(db tutorme.DB, auth *tutorme.Auth, clientID string) (*tutorme.Auth, error) {
	var createdAuth tutorme.Auth
	row := db.QueryRowx(insertToken, auth.Token, auth.AuthType, clientID)
	err := row.StructScan(&createdAuth)

	return &createdAuth, err
}

const blockClientQuery string = `
UPDATE auth
SET blocked = $1
WHERE client_id = $2 
`

func (au *AuthStore) BlockClient(db tutorme.DB, clientID string, blocked bool) error {
	_, err := db.Queryx(blockClientQuery, blocked, clientID)

	return err
}
