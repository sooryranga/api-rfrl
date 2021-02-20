package auth

import (
	"github.com/Arun4rangan/api-tutorme/src/client"
	"github.com/Arun4rangan/api-tutorme/src/db"
)

const (
	getByToken string = `
SELECT client.* FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.token =$1 AND auth.auth_type =$2 
LIMIT 1
	`
	getByEmail string = `
SELECT password_hash, client.* FROM auth 
JOIN client ON auth.client_id = client.id
WHERE auth.email =$1 AND auth.auth_type =$2 
LIMIT 1
	`
	insertEmailAuth string = `
INSERT INTO auth (email, password_hash, auth_type, client_id)
VALUES ($1, $2, $3, $4)
RETURNING id
	`
	insertToken string = `
INSERT INTO auth (token, auth_type, client_id) 
VALUES ($1, $2, $3) 
RETURNING id
	`
)

// GetByToken queries the database for token auth from providers
func GetByToken(db db.DB, token string, authType string) (*client.Client, error) {
	var c client.Client
	err := db.QueryRowx(getByToken, token, authType).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetByEmail queries the database for email auth
func GetByEmail(db db.DB, email string) (*client.Client, []byte, error) {

	type getByEmailStruct struct {
		client.Client
		PasswordHash []byte `db:"password_hash"`
	}
	var result getByEmailStruct

	row := db.QueryRowx(getByEmail, email, EMAIL)

	err := row.StructScan(&result)

	return &result.Client, result.PasswordHash, err
}

// CreateWithEmail creates auth row with email in db
func CreateWithEmail(db db.DB, auth *Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertEmailAuth, auth.Email, auth.PasswordHash, EMAIL, clientID)
	var id int = -1

	err := row.Scan(&id)

	return id, err
}

// CreateWithToken creates auth row with token in db
func CreateWithToken(db db.DB, auth *Auth, clientID string) (int, error) {
	row := db.QueryRowx(insertToken, auth.Token, auth.AuthType, clientID)
	var id int = -1
	err := row.Scan(&id)

	return id, err
}
