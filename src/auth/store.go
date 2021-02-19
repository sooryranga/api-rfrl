package auth

import (
	"github.com/Arun4rangan/api-tutorme/src/db"
)

const (
	getByToken string = `
SELECT client.*, password FROM auth 
JOIN client ON auth.id == client.id
WHERE auth.token =$1 AND auth.auth_type =$2 
LIMIT 1
	`
	getByEmail string = `
SELECT * FROM auth 
JOIN client ON auth.id == client.id
WHERE auth.email =$1 AND auth.auth_type =$2 
LIMIT 1
	`
	insertEmailAuth string = `
INSERT INTO auth (email, password_hash, auth_type)
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
func GetByToken(db db.DB, token string, authType string) (*Auth, error) {
	var m Auth
	err := db.QueryRowx(getByToken, token, authType).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByEmail queries the database for email auth
func GetByEmail(db db.DB, email string) (*Auth, error) {
	var m Auth
	err := db.QueryRowx(getByEmail, email, EMAIL).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateWithEmail creates auth row with email in db
func CreateWithEmail(db db.DB, auth *Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertEmailAuth, auth.Email, auth.PasswordHash, EMAIL)
	var id int = -1

	err := row.Scan(&id)

	return id, err
}

// CreateWithToken creates auth row with token in db
func CreateWithToken(db db.DB, auth *Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertToken, auth.Token, auth.AuthType)
	var id int = -1
	err := row.Scan(&id)

	return id, err
}
