package usecases

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// AuthUseCase holds all business related functions for auth
type AuthUseCase struct {
	db          *sqlx.DB
	authStore   rfrl.AuthStore
	clientStore rfrl.ClientStore
	fireStore   rfrl.FireStoreClient
}

// NewAuthUseCase creates new AuthUseCase
func NewAuthUseCase(
	db sqlx.DB,
	authStore rfrl.AuthStore,
	clientStore rfrl.ClientStore,
	fireStore rfrl.FireStoreClient,
) *AuthUseCase {
	return &AuthUseCase{&db, authStore, clientStore, fireStore}
}

// SignupWithToken allows user to sign up with token from google or linkedin auth
func (au *AuthUseCase) SignupWithToken(newClient *rfrl.Client, auth *rfrl.Auth) (*rfrl.Client, *rfrl.Auth, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = au.db.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	if *err != nil {
		return nil, nil, *err
	}

	var createdClient *rfrl.Client

	createdClient, *err = au.clientStore.CreateClient(tx, newClient)

	if *err != nil {
		return nil, nil, *err
	}

	var createdAuth *rfrl.Auth
	createdAuth, *err = au.authStore.CreateWithToken(tx, auth, createdClient.ID)

	if *err != nil {
		return nil, nil, *err
	}

	*err = au.fireStore.CreateClient(
		createdClient.ID,
		createdClient.Photo.String,
		createdClient.FirstName.String,
		createdClient.LastName.String,
	)

	if *err != nil {
		return nil, nil, *err
	}

	var firebaseToken string
	firebaseToken, *err = au.fireStore.CreateLoginToken(createdClient.ID)
	createdAuth.FirebaseToken = firebaseToken

	return createdClient, createdAuth, *err
}

// SignupGoogle allows user to sign up with google
func (au *AuthUseCase) SignupGoogle(
	token string,
	email string,
	firstName string,
	lastName string,
	photo string,
	about string,
	isTutor null.Bool,
) (*rfrl.Client, *rfrl.Auth, error) {

	newClient := rfrl.NewClient(firstName, lastName, about, email, photo, isTutor, "", "", null.Int{}, "")
	auth := rfrl.Auth{
		AuthType: null.NewString(rfrl.GOOGLE, true),
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
	isTutor null.Bool,
) (*rfrl.Client, *rfrl.Auth, error) {

	newClient := rfrl.NewClient(firstName, lastName, about, email, photo, isTutor, "", "", null.Int{}, "")

	auth := rfrl.Auth{
		AuthType: null.NewString(rfrl.LINKEDIN, true),
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
	isTutor null.Bool,
) (*rfrl.Client, *rfrl.Auth, error) {
	hash, hashError := hashAndSalt([]byte(password))

	if hashError != nil {
		return nil, nil, hashError
	}

	newClient := rfrl.NewClient(firstName, lastName, about, email, photo, isTutor, "", "", null.Int{}, "")
	auth := rfrl.Auth{
		Email:        null.StringFrom(email),
		PasswordHash: hash,
	}

	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = au.db.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	if *err != nil {
		return nil, nil, *err
	}

	var createdClient *rfrl.Client

	createdClient, *err = au.clientStore.CreateClient(tx, newClient)

	if *err != nil {
		return nil, nil, *err
	}

	var createdAuth *rfrl.Auth
	createdAuth, *err = au.authStore.CreateWithEmail(tx, &auth, createdClient.ID)

	if *err != nil {
		return nil, nil, *err
	}

	*err = au.fireStore.CreateClient(
		createdClient.ID,
		createdClient.Photo.String,
		createdClient.FirstName.String,
		createdClient.LastName.String,
	)

	if *err != nil {
		return nil, nil, *err
	}

	var firebaseToken string
	firebaseToken, *err = au.fireStore.CreateLoginToken(createdClient.ID)
	createdAuth.FirebaseToken = firebaseToken

	return createdClient, createdAuth, *err
}

// LoginEmail allows user to login with email by checking password hash against the has the passed in
func (au *AuthUseCase) LoginEmail(email string, password string) (*rfrl.Client, *rfrl.Auth, error) {
	c, auth, err := au.authStore.GetByEmail(au.db, email)

	if err != nil {
		return nil, nil, errors.Wrap(err, "GetByEmail")
	}

	err = bcrypt.CompareHashAndPassword(
		auth.PasswordHash,
		[]byte(password),
	)

	if err != nil {
		return nil, nil, errors.Wrap(err, "CompareHashAndPassword")
	}

	firebaseToken, err := au.fireStore.CreateLoginToken(c.ID)
	auth.FirebaseToken = firebaseToken

	return c, auth, err
}

func (au *AuthUseCase) LoginWithJWT(clientID string) (*rfrl.Client, *rfrl.Auth, error) {
	c, err := au.clientStore.GetClientFromID(au.db, clientID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.Wrap(err, fmt.Sprintf("Client not found with %s", clientID))
		}
		return nil, nil, errors.Wrap(err, "LoginWithJWT")
	}

	auth, err := au.authStore.GetByClientID(au.db, clientID)

	if err != nil {
		return nil, nil, err
	}

	firebaseToken, err := au.fireStore.CreateLoginToken(c.ID)
	auth.FirebaseToken = firebaseToken

	return c, auth, err
}

// LoginGoogle allow user to login with their google oauth token
func (au *AuthUseCase) LoginGoogle(token string) (*rfrl.Client, *rfrl.Auth, error) {
	cl, a, err := au.authStore.GetByToken(au.db, token, rfrl.GOOGLE)

	if err != nil {
		return cl, a, err
	}

	a.FirebaseToken, err = au.fireStore.CreateLoginToken(cl.ID)
	return cl, a, err
}

// LoginLinkedIn allow user to login with their linkedin
func (au *AuthUseCase) LoginLinkedIn(token string) (*rfrl.Client, *rfrl.Auth, error) {
	cl, a, err := au.authStore.GetByToken(au.db, token, rfrl.LINKEDIN)

	if err != nil {
		return cl, a, err
	}

	a.FirebaseToken, err = au.fireStore.CreateLoginToken(cl.ID)
	return cl, a, err
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
func (au *AuthUseCase) GenerateToken(claims *rfrl.JWTClaims, signingKey *rsa.PrivateKey) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (au *AuthUseCase) UpdateSignUpFlow(clientID string, stage rfrl.SignUpFlow) error {
	return au.authStore.UpdateSignUpFlow(au.db, clientID, stage)
}

func (au *AuthUseCase) BlockClient(clientID string, blocked bool) error {
	return au.authStore.BlockClient(au.db, clientID, blocked)
}
