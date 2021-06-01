package tutorme

import (
	"crypto/rsa"
	"database/sql/driver"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
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

type SignUpFlow int

const (
	TypeOfUser        = 0
	BasicInfo         = 1
	CompanyEmail      = 2
	RegisterDocuments = 3
	DoneSignUp        = 100
)

var signUpFlowValueToJSON = map[int]string{
	0:   `"TypeOfUser"`,
	1:   `"BasicInfo"`,
	2:   `"CompanyEmail"`,
	3:   `"RegisterDocuments"`,
	100: `"DoneSignUp"`,
}

var signUpFlowReadableToValue = map[string]int{
	"TypeOfUser":        0,
	"BasicInfo":         1,
	"CompanyEmail":      2,
	"RegisterDocuments": 3,
	"DoneSignUp":        100,
}

func (d SignUpFlow) MarshalJSON() ([]byte, error) {
	return []byte(signUpFlowValueToJSON[int(d)]), nil
}

func (d *SignUpFlow) UnmarshalJSON(b []byte) error {
	var signUpFlowString string
	err := json.Unmarshal(b, &signUpFlowString)

	if err != nil {
		return err
	}

	n, ok := signUpFlowReadableToValue[signUpFlowString]
	if !ok {
		return errors.Errorf("Invalid for SignUpFlow (%s)", string(b))
	}
	*d = SignUpFlow(n)

	return nil
}

func (s *SignUpFlow) Scan(src interface{}) error {
	var value int64
	switch src.(type) {
	case int64:
		value = src.(int64)
	default:
		errors.New("Invalid type for SignUpFlow")
	}
	*s = SignUpFlow(value)
	return nil
}

func (s SignUpFlow) Value() (driver.Value, error) {
	values := map[SignUpFlow]interface{}{
		TypeOfUser:        nil,
		BasicInfo:         nil,
		CompanyEmail:      nil,
		DoneSignUp:        nil,
		RegisterDocuments: nil,
	}

	if _, ok := values[s]; !ok {
		return nil, errors.New("Wrong value for CustomType")
	}

	return int(s), nil
}

const AlgorithmRS256 string = "RS256"

func GetClaims(c echo.Context) (*JWTClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || !token.Valid {
		return nil, errors.New("Token is not valid")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Claims is not valid")
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
	ID            int         `db:"id" json:"-"`
	CreatedAt     time.Time   `db:"created_at" json:"-"`
	UpdatedAt     time.Time   `db:"updated_at" json:"-"`
	Token         null.String `db:"token" json:"-"`
	AuthType      null.String `db:"auth_type" json:"-"`
	Email         null.String `db:"email" json:"-"`
	PasswordHash  []byte      `db:"password_hash" json:"-"`
	ClientID      string      `db:"client_id" json:"-"`
	Stage         SignUpFlow  `db:"sign_up_flow" json:"signUpStage"`
	Blocked       bool        `db:"blocked" json:"-"`
	FirebaseToken string      `json:"firebaseToken"`
}

type AuthStore interface {
	GetByClientID(db DB, clientID string) (*Auth, error)
	GetByToken(db DB, token string, authType string) (*Client, *Auth, error)
	GetByEmail(db DB, email string) (*Client, *Auth, error)
	CreateWithEmail(db DB, auth *Auth, clientID string) (*Auth, error)
	CreateWithToken(db DB, auth *Auth, clientID string) (*Auth, error)
	CheckEmailAuthExists(db DB, clientID string, email string) (bool, error)
	UpdateAuthEmail(db DB, clientID string, email string) error
	UpdateSignUpFlow(db DB, clientID string, stage SignUpFlow) error
	BlockClient(db DB, clientID string, blocked bool) error
}

type AuthUseCase interface {
	SignupWithToken(newClient *Client, auth *Auth) (*Client, *Auth, error)
	SignupGoogle(token string, email string, firstName string, lastName string, photo string, about string, isTutor null.Bool) (*Client, *Auth, error)
	SignupLinkedIn(token string, email string, firstName string, lastName string, photo string, about string, isTutor null.Bool) (*Client, *Auth, error)
	SignupEmail(password string, token string, email string, firstName string, lastName string, photo string, about string, isTutor null.Bool) (*Client, *Auth, error)
	LoginEmail(email string, password string) (*Client, *Auth, error)
	LoginGoogle(token string) (*Client, *Auth, error)
	LoginLinkedIn(token string) (*Client, *Auth, error)
	LoginWithJWT(clientID string) (*Client, *Auth, error)
	GenerateToken(claims *JWTClaims, signingKey *rsa.PrivateKey) (string, error)
	UpdateSignUpFlow(clientID string, stage SignUpFlow) error
	BlockClient(clientID string, blocked bool) error
}
