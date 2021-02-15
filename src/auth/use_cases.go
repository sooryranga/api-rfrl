package auth

import (
	"database/sql"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) signupGoogle(token string) (string, error) {
	auth := Auth{
		AuthType: GOOGLE,
		Token:    sql.NullString{String: token, Valid: true},
	}
	_, err := h.authStore.CreateWithToken(&auth)

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
	_, err := h.authStore.CreateWithToken(&auth)

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
	_, err = h.authStore.CreateWithEmail(&auth)

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
	auth, err := h.authStore.GetByEmail(email)

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
	_, err := h.authStore.GetByToken(token, GOOGLE)

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
	_, err := h.authStore.GetByToken(token, LINKEDIN)

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
