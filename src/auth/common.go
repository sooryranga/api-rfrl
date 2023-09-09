package auth

import (
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

type Handler struct {
	authStore Store
}

func NewHandler(au Store) *Handler {
	return &Handler{
		authStore: au,
	}
}

// Types definition
const (
	GOOGLE   = "google"
	LINKEDIN = "linkedin"
	EMAIL    = "email"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// GenerateToken creates token
func GenerateToken(claims *jwtCustomClaims) (string, error) {
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
