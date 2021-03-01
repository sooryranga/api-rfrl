package tutorme

import (
	"crypto/rsa"
	"database/sql"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

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

const AlgorithmRS256 string = "rsa"

func GetClaims(c echo.Context) *JWTClaims {
	user := c.Get("user").(*jwt.Token)
	return user.Claims.(*JWTClaims)
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

// Auth model
type Auth struct {
	ID           string         `db:"id"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	Token        sql.NullString `db:"token"`
	AuthType     string         `db:"auth_type"`
	Email        sql.NullString `db:"email"`
	PasswordHash []byte         `db:"password_hash"`
	ClientID     string         `db:"client_id"`
}

type AuthStore interface {
	GetByToken(db DB, token string, authType string) (*Client, error)
	GetByEmail(db DB, email string) (*Client, []byte, error)
	CreateWithEmail(db DB, auth *Auth, clientID string) (int, error)
	CreateWithToken(db DB, auth *Auth, clientID string) (int, error)
}

type AuthUseCase interface {
	SignupWithToken(newClient *Client, auth *Auth) (*Client, error)
	SignupGoogle(token string, email string, firstName string, lastName string, photo string, about string) (*Client, error)
	SignupLinkedIn(token string, email string, firstName string, lastName string, photo string, about string) (*Client, error)
	SignupEmail(password string, token string, email string, firstName string, lastName string, photo string, about string) (*Client, error)
	LoginEmail(email string, password string) (*Client, error)
	LoginGoogle(token string) (*Client, error)
	LoginLinkedIn(token string) (*Client, error)
	GenerateToken(claims *JWTClaims, signingKey *rsa.PrivateKey) (string, error)
}
