package usecases

import (
	"crypto/rsa"
	"fmt"
	"log"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// AuthUseCase holds all business related functions for auth
type AuthUseCase struct {
	db          *sqlx.DB
	authStore   tutorme.AuthStore
	clientStore tutorme.ClientStore
}

// NewAuthUseCase creates new AuthUseCase
func NewAuthUseCase(db sqlx.DB, authStore tutorme.AuthStore, clientStore tutorme.ClientStore) *AuthUseCase {
	return &AuthUseCase{&db, authStore, clientStore}
}

// SignupWithToken allows user to sign up with token from google or linkedin auth
func (au *AuthUseCase) SignupWithToken(newClient *tutorme.Client, auth *tutorme.Auth) (*tutorme.Client, error) {
	tx, err := au.db.Beginx()

	if err != nil {
		return nil, err
	}

	newClient, err = au.clientStore.CreateClient(tx, newClient)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	_, err = au.authStore.CreateWithToken(tx, auth, newClient.ID)

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

// SignupGoogle allows user to sign up with google
func (au *AuthUseCase) SignupGoogle(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*tutorme.Client, error) {

	newClient := tutorme.NewClient(firstName, lastName, about, email, photo)
	auth := tutorme.Auth{
		AuthType: tutorme.GOOGLE,
		Token:    null.StringFrom(token),
	}

	return au.SignupWithToken(newClient, &auth)
}

// SignupLinkedIn allows user to sign up with linkedin
func (au *AuthUseCase) SignupLinkedIn(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*tutorme.Client, error) {

	newClient := tutorme.NewClient(firstName, lastName, about, email, photo)

	auth := tutorme.Auth{
		AuthType: tutorme.LINKEDIN,
		Token:    null.StringFrom(token),
	}

	return au.SignupWithToken(newClient, &auth)
}

// SignupEmail allows user to signup with email
func (au *AuthUseCase) SignupEmail(
	password string,
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
) (*tutorme.Client, error) {
	hash, err := hashAndSalt([]byte(password))

	if err != nil {
		return nil, err
	}

	newClient := tutorme.NewClient(firstName, lastName, about, email, photo)
	auth := tutorme.Auth{
		Email:        null.StringFrom(email),
		PasswordHash: hash,
	}

	tx, err := au.db.Beginx()
	if err != nil {
		return nil, err
	}

	newClient, err = au.clientStore.CreateClient(tx, newClient)
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	_, err = au.authStore.CreateWithEmail(tx, &auth, newClient.ID)

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

// LoginEmail allows user to login with email by checking password hash against the has the passed in
func (au *AuthUseCase) LoginEmail(email string, password string) (*tutorme.Client, error) {
	c, passwordHash, err := au.authStore.GetByEmail(au.db, email)

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

func (au *AuthUseCase) LoginWithJWT(clientID string) (*tutorme.Client, error) {
	c, err := au.clientStore.GetClientFromID(au.db, clientID)

	if err != nil {
		return nil, errors.Wrap(err, "LoginWithJWT")
	}

	return c, nil
}

// LoginGoogle allow user to login with their google oauth token
func (au *AuthUseCase) LoginGoogle(token string) (*tutorme.Client, error) {
	return au.authStore.GetByToken(au.db, token, tutorme.GOOGLE)
}

// LoginLinkedIn allow user to login with their linkedin
func (au *AuthUseCase) LoginLinkedIn(token string) (*tutorme.Client, error) {
	return au.authStore.GetByToken(au.db, token, tutorme.LINKEDIN)
}

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

// GenerateToken creates token
func (au *AuthUseCase) GenerateToken(claims *tutorme.JWTClaims, signingKey *rsa.PrivateKey) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return t, nil
}
