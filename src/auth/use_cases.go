package auth

import (
	"database/sql"
	"time"

	"github.com/Arun4rangan/api-tutorme/src/user"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) signupGoogle(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (string, error) {

	tx, err := h.db.Begin()

	if err != nil {
		return "", err
	}
	users := user.User{
		firstName: sql.NullString{String: firstName},
		LastName:  lastName,
	}
	auth := Auth{
		AuthType: GOOGLE,
		Token:    token,
	}
	_, err := CreateWithToken(h.db, &auth, nil)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		"",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}

func (h *Handler) signupLinkedIn(token string) (string, error) {
	auth := Auth{
		AuthType: LINKEDIN,
		Token:    sql.NullString{String: token, Valid: true},
	}
	_, err := CreateWithToken(h.db, &auth)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		"",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}

func (h *Handler) signupEmail(email string, password string) (string, error) {
	hash, err := hashAndSalt([]byte(password))

	if err != nil {
		return "", err
	}

	auth := Auth{
		Email:        sql.NullString{String: email, Valid: true},
		PasswordHash: hash,
	}
	_, err = CreateWithEmail(h.db, &auth)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}

func (h *Handler) loginEmail(email string, password string) (string, error) {
	auth, err := GetByEmail(h.db, email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		auth.PasswordHash,
		[]byte(password),
	)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}

func (h *Handler) loginGoogle(token string) (string, error) {
	_, err := GetByToken(h.db, token, GOOGLE)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		"",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}

func (h *Handler) loginLinkedIn(token string) (string, error) {
	_, err := GetByToken(h.db, token, LINKEDIN)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		"",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return h.GenerateToken(claims)
}
