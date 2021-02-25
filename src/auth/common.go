package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles all request within auth
type Handler struct {
	db         *sqlx.DB
	signingKey *rsa.PrivateKey
}

// NewHandler creates a handler
func NewHandler(db *sqlx.DB, key *rsa.PrivateKey) *Handler {
	return &Handler{
		db:         db,
		signingKey: key,
	}
}

// JWTClaims are custom claims extending default ones.
type JWTClaims struct {
	ClientID string `json:"id"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
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
	publicRsaFileLoc := os.Getenv("PUBLIC_RSA_FILE")
	keyData, err := ioutil.ReadFile(publicRsaFileLoc)

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// GetSigningKey generate private key from id_rsa
func GetSigningKey() (*rsa.PrivateKey, error) {
	rsaFileLoc := os.Getenv("RSA_FILE")
	keyData, err := ioutil.ReadFile(rsaFileLoc)

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

func GetClaims(c echo.Context) *JWTClaims {
	user := c.Get("user").(*jwt.Token)
	return user.Claims.(*JWTClaims)
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
