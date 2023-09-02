package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// Login is used to login an user
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	t, err := generateToken(claims)

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
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
