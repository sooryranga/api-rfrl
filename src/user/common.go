package user

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
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

// Types definition
const (
	GOOGLE   = "google"
	LINKEDIN = "linkedin"
	EMAIL    = "email"
)

// JWTClaims are custom claims extending default ones.
type JWTClaims struct {
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// GenerateToken creates token
func GenerateToken(claims *JWTClaims) (string, error) {
	// TODO:read id_rsa once
	keyData, err := ioutil.ReadFile("./id_rsa")

	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		return "", err
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}
