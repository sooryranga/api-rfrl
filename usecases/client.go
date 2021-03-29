package usecases

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

// ClientUseCase holds all business related functions for client
type ClientUseCase struct {
	db          *sqlx.DB
	clientStore tutorme.ClientStore
}

// NewClientUseCase creates new ClientUseCase
func NewClientUseCase(db sqlx.DB, clientStore tutorme.ClientStore) *ClientUseCase {
	return &ClientUseCase{&db, clientStore}
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
	return cl.clientStore.CreateClient(cl.db, client)
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
	client := tutorme.NewClient(
		firstName,
		lastName,
		about,
		email,
		photo,
		isTutor,
	)

	return cl.clientStore.UpdateClient(cl.db, id, client)
}

// GetClient use case to get existing client
func (cl *ClientUseCase) GetClient(id string) (*tutorme.Client, error) {
	return cl.clientStore.GetClientFromID(cl.db, id)
}

func (cl *ClientUseCase) GetClients(options tutorme.GetClientsOptions) (*[]tutorme.Client, error) {
	return cl.clientStore.GetClients(cl.db, options)
}
