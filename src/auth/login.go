package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

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

type (
	// LoginPayload is the struct used to hold payload from /login
	LoginPayload struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Token    string `json:"token"`
		Password string `json:"password" validate:"omitempty,gte=10"`
		Type     string `json:"type" validate:"required,oneof= google linkedin email"`
	}
)

// LoginPayloadValidation validates user inputs
func LoginPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(LoginPayload)

	switch payload.Type {
	case GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

func loginEmail(c echo.Context) (string, error) {
	email := "arun.ranga@hotmail.ca"

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return generateToken(claims)
}

func loginGoogle(c echo.Context) (string, error) {
	email := "arun.ranga@hotmail.ca"

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return generateToken(claims)
}

func loginLinkedIn(c echo.Context) (string, error) {
	email := "arun.ranga@hotmail.ca"

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return generateToken(claims)
}

// Login is used to login an user
func Login(c echo.Context) error {
	payload := LoginPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var token string
	var error error

	switch payload.Type {
	case GOOGLE:
		token, error = loginGoogle(c)
	case LINKEDIN:
		token, error = loginLinkedIn(c)
	case EMAIL:
		token, error = loginEmail(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if error != nil {
		panic(fmt.Sprintf("%v", error))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

// GenerateToken creates token
func generateToken(claims *jwtCustomClaims) (string, error) {
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
