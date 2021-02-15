package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles all request within auth
type Handler struct {
	authStore  Store
	signingKey *rsa.PrivateKey
}

// NewHandler creates a handler
func NewHandler(au Store, key *rsa.PrivateKey) *Handler {
	return &Handler{
		authStore:  au,
		signingKey: key,
	}
}

// JWTClaims are custom claims extending default ones.
type JWTClaims struct {
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// Types definition
const (
	GOOGLE   = "google"
	LINKEDIN = "linkedin"
	EMAIL    = "email"
)

func hashAndSalt(pwd []byte) ([]byte, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return hash, nil
}

// GetVerifyingKey generate public key from id_rsa.pub
func GetVerifyingKey() (*rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile("./id_rsa.pub")

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// GetSigningKey generate private key from id_rsa
func GetSigningKey() (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile("./id_rsa")

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

// GenerateToken creates token
func (h *Handler) GenerateToken(claims *JWTClaims) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(h.signingKey)
	if err != nil {
		return "", err
	}

	return t, nil
}
