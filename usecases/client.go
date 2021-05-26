package usecases

import (
	"database/sql"
	"errors"
	"strings"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

// ClientUseCase holds all business related functions for client
type ClientUseCase struct {
	db           *sqlx.DB
	clientStore  tutorme.ClientStore
	emailer      tutorme.EmailerUseCase
	authStore    tutorme.AuthStore
	fireStore    tutorme.FireStoreClient
	companyStore tutorme.CompanyStore
}

// NewClientUseCase creates new ClientUseCase
func NewClientUseCase(
	db sqlx.DB,
	clientStore tutorme.ClientStore,
	authStore tutorme.AuthStore,
	emailer tutorme.EmailerUseCase,
	fireStore tutorme.FireStoreClient,
	companyStore tutorme.CompanyStore,
) *ClientUseCase {
	return &ClientUseCase{
		&db,
		clientStore,
		emailer,
		authStore,
		fireStore,
		companyStore,
	}
}

// CreateClient use case to create a new client
func (cl *ClientUseCase) CreateClient(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
	isTutor null.Bool,
) (*tutorme.Client, error) {
	client := tutorme.NewClient(
		firstName,
		lastName,
		about,
		email,
		photo,
		isTutor,
	)
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = cl.db.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	if *err != nil {
		return nil, *err
	}

	var createdClient *tutorme.Client
	createdClient, *err = cl.clientStore.CreateClient(cl.db, client)

	if *err != nil {
		return nil, *err
	}

	*err = cl.fireStore.CreateClient(
		createdClient.ID,
		createdClient.Photo.String,
		createdClient.FirstName.String,
		createdClient.LastName.String,
	)

	return createdClient, *err
}

// UpdateClient use case to update a new client
func (cl *ClientUseCase) UpdateClient(
	id string,
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
	isTutor null.Bool,
) (*tutorme.Client, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = cl.db.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	client := tutorme.NewClient(
		firstName,
		lastName,
		about,
		email,
		photo,
		isTutor,
	)

	var updatedClient *tutorme.Client
	updatedClient, *err = cl.clientStore.UpdateClient(cl.db, id, client)

	if *err != nil {
		return nil, *err
	}

	if client.FirstName.Valid || client.LastName.Valid {
		*err = cl.fireStore.UpdateClient(id, client.Photo, updatedClient.FirstName, updatedClient.LastName)
	} else {
		*err = cl.fireStore.UpdateClient(id, client.Photo, client.FirstName, client.LastName)
	}

	return updatedClient, *err
}

// GetClient use case to get existing client
func (cl *ClientUseCase) GetClient(id string) (*tutorme.Client, error) {
	return cl.clientStore.GetClientFromID(cl.db, id)
}

func (cl *ClientUseCase) GetClients(options tutorme.GetClientsOptions) (*[]tutorme.Client, error) {
	return cl.clientStore.GetClients(cl.db, options)
}

func (cl *ClientUseCase) CreateEmailVerification(clientID string, email string, emailType string) error {

	passcode, err := cl.emailer.SendEmailVerification(email)

	if err != nil {
		return err
	}

	if emailType == tutorme.UserEmail {
		exists, err := cl.authStore.CheckEmailAuthExists(cl.db, clientID, email)

		if err != nil {
			return err
		}

		if exists {
			return errors.New("Email is already in use")
		}
	}

	return cl.clientStore.CreateEmailVerification(cl.db, clientID, email, emailType, passcode)
}

func (cl *ClientUseCase) verifyUserEmail(tx *sqlx.Tx, clientID string, email string) (*tutorme.Client, error) {
	err := cl.authStore.UpdateAuthEmail(tx, clientID, email)

	if err != nil {
		return nil, err
	}

	client := tutorme.Client{Email: null.NewString(email, true), VerifiedEmail: null.NewBool(true, true)}

	updatedClient, err := cl.clientStore.UpdateClient(tx, clientID, &client)

	return updatedClient, err
}

func (cl *ClientUseCase) verifyWorkEmail(tx *sqlx.Tx, clientID string, email string) (*tutorme.Client, error) {
	at := strings.LastIndex(email, "@")
	_, domain := email[:at], email[at+1:]

	companyID, err := cl.companyStore.GetCompanyIDFromEmailDomain(tx, domain)

	if err == sql.ErrNoRows {
		err = cl.companyStore.CreateCompanyEmailDomain(tx, domain)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	client := tutorme.Client{
		WorkEmail:         null.NewString(email, true),
		VerifiedWorkEmail: null.NewBool(true, true),
		CompanyID:         companyID,
	}

	updatedClient, err := cl.clientStore.UpdateClient(tx, clientID, &client)

	return updatedClient, err
}

func (cl *ClientUseCase) VerifyEmail(clientID string, email string, emailType string, passCode string) (*tutorme.Client, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = cl.db.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	*err = cl.clientStore.VerifyEmail(tx, clientID, email, emailType, passCode)

	if *err != nil {
		return nil, *err
	}

	var updatedClient *tutorme.Client
	if emailType == tutorme.WorkEmail {
		updatedClient, *err = cl.verifyWorkEmail(tx, clientID, email)
	} else if emailType == tutorme.UserEmail {
		updatedClient, *err = cl.verifyUserEmail(tx, clientID, email)
	}

	if *err != nil {
		return nil, *err
	}

	return updatedClient, nil
}

func (cl *ClientUseCase) GetVerificationEmail(clientID string, emailType string) (string, error) {
	return cl.clientStore.GetVerificationEmail(cl.db, clientID, emailType)
}

func (cl *ClientUseCase) DeleteVerificationEmail(clientID string, emailType string) error {
	return cl.clientStore.DeleteVerificationEmail(cl.db, clientID, emailType)
}

func (cl *ClientUseCase) GetClientEvents(clientID string, start null.Time, end null.Time, state null.String) (*[]tutorme.Event, error) {
	return cl.clientStore.GetRelatedEventsByClientIDs(cl.db, []string{clientID}, start, end, state)
}

func (cl *ClientUseCase) CreateOrUpdateClientEducation(clientID string, institution string, degree string, fieldOfStudy string, startYear int, endYear int) error {
	education := tutorme.NewEducation(institution, degree, fieldOfStudy, startYear, endYear)
	err := cl.clientStore.CreateOrUpdateClientEducation(cl.db, clientID, education)

	return err
}
