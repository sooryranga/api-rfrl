package store

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// ClientStore holds all store related functions for client
type ClientStore struct{}

// NewClientStore creates new clientStore
func NewClientStore() *ClientStore {
	return &ClientStore{}
}

const (
	getClientByID string = `
SELECT * FROM client
WHERE client.id = $1
	`
	insertClient string = `
INSERT INTO client (first_name, last_name, about, email, photo)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
	`
)

// GetClientFromID queries the database for client with id
func (cl *ClientStore) GetClientFromID(db tutorme.DB, id string) (*tutorme.Client, error) {
	var m tutorme.Client
	err := db.QueryRowx(getClientByID, id).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateClient creates a new row for a client in the database
func (cl *ClientStore) CreateClient(db tutorme.DB, client *tutorme.Client) (*tutorme.Client, error) {
	row := db.QueryRowx(
		insertClient,
		client.FirstName,
		client.LastName,
		client.About,
		client.Email,
		client.Photo,
	)

	var m tutorme.Client

	err := row.StructScan(&m)
	return &m, errors.Wrap(err, "CreateClient")
}

// UpdateClient updates a client in the database
func (cl *ClientStore) UpdateClient(db tutorme.DB, ID string, client *tutorme.Client) (*tutorme.Client, error) {
	query := sq.Update("client")
	if client.FirstName.Valid {
		query = query.Set("first_name", client.FirstName)
	}
	if client.LastName.Valid {
		query = query.Set("last_name", client.LastName)
	}
	if client.About.Valid {
		query = query.Set("about", client.About)
	}
	if client.Photo.Valid {
		query = query.Set("photo", client.Photo)
	}
	if client.Email.Valid {
		query = query.Set("email", client.Email)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": ID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(
		sql,
		args...,
	)

	var m tutorme.Client

	err = row.StructScan(&m)
	return &m, err
}
