package auth

import (
	"database/sql"
	"time"

	"github.com/Arun4rangan/api-tutorme/src/client"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) signupWithToken(newClient *client.Client, auth *Auth) (*client.Client, error) {
	tx, err := h.db.Beginx()

	if err != nil {
		return nil, err
	}

	newClient, err = client.CreateClient(tx, newClient)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	_, err = CreateWithToken(h.db, auth, newClient.ID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}
	return newClient, nil
}

func (h *Handler) signupGoogle(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*client.Client, error) {

	newClient := client.NewClient(firstName, lastName, about, email, photo)
	auth := Auth{
		AuthType: GOOGLE,
		Token:    sql.NullString{String: token, Valid: true},
	}

	return h.signupWithToken(newClient, &auth)
}

func (h *Handler) signupLinkedIn(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*client.Client, error) {

	newClient := client.NewClient(firstName, lastName, about, email, photo)

	auth := Auth{
		AuthType: LINKEDIN,
		Token:    sql.NullString{String: token, Valid: true},
	}

	return h.signupWithToken(newClient, &auth)
}

func (h *Handler) signupEmail(
	password string,
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*client.Client, error) {
	hash, err := hashAndSalt([]byte(password))

	if err != nil {
		return "", err
	}

	newClient := client.NewClient(firstName, lastName, about, email, photo)
	auth := Auth{
		Email:        sql.NullString{String: email, Valid: true},
		PasswordHash: hash,
	}

	tx, err := h.db.Beginx()
	if err != nil {
		return nil, err
	}

	newClient, err = client.CreateClient(tx, newClient)
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	_, err = CreateWithEmail(tx, &auth, newClient.ID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}
	return newClient, nil
}

func (h *Handler) loginEmail(email string, password string) (*client.Client, error) {
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

	return client.GetClientFromID(h.db, auth.ClientID)
}

func (h *Handler) loginGoogle(token string) (string, error) {
	auth, err := GetByToken(h.db, token, GOOGLE)

	if err != nil {
		return "", err
	}

	return "", err
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
