package client

import (
	"github.com/Arun4rangan/api-tutorme/src/db"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

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
func GetClientFromID(db db.DB, id string) (*Client, error) {
	var m Client
	err := db.QueryRowx(getClientByID, id).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateClient creates a new row for a client in the database
func CreateClient(db db.DB, client *Client) (*Client, error) {
	row := db.QueryRowx(
		insertClient,
		client.FirstName,
		client.LastName,
		client.About,
		client.Email,
		client.Photo,
	)

	var m Client

	err := row.StructScan(&m)
	return &m, errors.Wrap(err, "CreateClient")
}

// UpdateClient updates a client in the database
func UpdateClient(db db.DB, ID string, client *Client) (*Client, error) {
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

	var m Client

	err = row.StructScan(&m)
	return &m, err
}
