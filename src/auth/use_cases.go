package auth

import (
	"database/sql"
	"fmt"

	"github.com/Arun4rangan/api-tutorme/src/client"
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

	_, err = CreateWithToken(tx, auth, newClient.ID)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, errors.Wrap(err, fmt.Sprintf("%v", newClient.ID))
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
		return nil, err
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
	c, passwordHash, err := GetByEmail(h.db, email)

	if err != nil {
		return nil, errors.Wrap(err, "GetByEmail")
	}

	err = bcrypt.CompareHashAndPassword(
		passwordHash,
		[]byte(password),
	)

	if err != nil {
		return nil, errors.Wrap(err, "CompareHashAndPassword")
	}

	return c, nil
}

func (h *Handler) loginGoogle(token string) (*client.Client, error) {
	return GetByToken(h.db, token, GOOGLE)
}

func (h *Handler) loginLinkedIn(token string) (*client.Client, error) {
	return GetByToken(h.db, token, LINKEDIN)
}
