package store

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

// ClientStore holds all store related functions for client
type ClientStore struct{}

// NewClientStore creates new clientStore
func NewClientStore() *ClientStore {
	return &ClientStore{}
}

const (
	getClientByIDSQL string = `
SELECT * FROM client
WHERE client.id = $1
	`
	getClientByIDsSQL string = `
	SELECT * FROM client
	WHERE client.id IN (?)
`
)

func (cl *ClientStore) GetClients(db tutorme.DB, options tutorme.GetClientsOptions) (*[]tutorme.Client, error) {
	query := sq.Select("*").From("client")

	if options.IsTutor.Valid {
		query = query.Where(sq.Eq{"is_tutor": options.IsTutor.Bool})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	clients := make([]tutorme.Client, 0)

	if err != nil {
		return &clients, err
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return &clients, err
	}

	for rows.Next() {
		var client tutorme.Client

		err := rows.StructScan(&client)

		if err != nil {
			return &clients, err
		}

		clients = append(clients, client)
	}

	return &clients, nil
}

// GetClientFromID queries the database for client with id
func (cl *ClientStore) GetClientFromID(db tutorme.DB, id string) (*tutorme.Client, error) {
	var m tutorme.Client
	err := db.QueryRowx(getClientByIDSQL, id).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetClientFromIDs queries the database for client with ids
func getClientFromIDs(db tutorme.DB, ids []string) (*[]tutorme.Client, error) {
	query, args, err := sqlx.In(getClientByIDsSQL, ids)

	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)

	clients := make([]tutorme.Client, 0)

	for rows.Next() {
		var c tutorme.Client
		err := rows.StructScan(&c)
		if err != nil {
			return nil, err
		}

		clients = append(clients, c)
	}

	return &clients, err
}

func (cl *ClientStore) GetClientFromIDs(db tutorme.DB, ids []string) (*[]tutorme.Client, error) {
	return getClientFromIDs(db, ids)
}

// CreateClient creates a new row for a client in the database
func (cl *ClientStore) CreateClient(db tutorme.DB, client *tutorme.Client) (*tutorme.Client, error) {
	columns := []string{"first_name", "last_name", "about", "email", "photo"}
	values := make([]interface{}, 0)
	values = append(values,
		client.FirstName,
		client.LastName,
		client.About,
		client.Email,
		client.Photo,
	)

	if client.IsTutor.Valid {
		columns = append(columns, "is_tutor")
		values = append(values, client.IsTutor)
	}

	query := sq.Insert("client").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING *")

	sql, args, err := query.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "CreateClient")
	}

	row := db.QueryRowx(
		sql,
		args...,
	)

	var m tutorme.Client

	err = row.StructScan(&m)
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
	if client.IsTutor.Valid {
		query = query.Set("is_tutor", client.IsTutor)
	}
	if client.WorkEmail.Valid {
		query = query.Set("work_email", client.WorkEmail)
	}
	if client.VerifiedWorkEmail.Valid {
		query = query.Set("verified_work_email", client.VerifiedWorkEmail)
	}
	if client.Email.Valid {
		query = query.Set("email", client.Email)
	}
	if client.VerifiedEmail.Valid {
		query = query.Set("verified_email", client.VerifiedEmail)
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

const createEmailVerification string = `
INSERT INTO email_verification(client_id, email, email_type, pass_code)
VALUES ($1, $2, $3, $4)
`

func (cl *ClientStore) CreateEmailVerification(
	db tutorme.DB,
	clientID string,
	email string,
	emailType string,
	passCode string,
) error {
	_, err := db.Queryx(createEmailVerification, clientID, email, emailType, passCode)

	return err
}

const verifyEmail string = `
	SELECT id 
	FROM email_verification
	WHERE client_id = $1
	AND email = $2
	AND email_type = $3
	AND pass_code = $4
	LIMIT 1
`

const deleteVerifyEmailFromId string = `
	DELETE FROM email_verification
	WHERE id = $1
`

func (cl *ClientStore) VerifyEmail(
	db tutorme.DB,
	clientID string,
	email string,
	emailType string,
	passCode string,
) error {
	var id null.Int

	err := db.QueryRowx(verifyEmail, clientID, email, emailType, passCode).Scan(&id)

	if err != nil {
		return err
	}

	if id.Valid == false {
		return errors.New("Could not find email to verify")
	}

	_, err = db.Queryx(deleteVerifyEmailFromId, id)

	return err
}

const getVerificationEmail string = `
SELECT email FROM email_verification
WHERE client_id = $1 AND email_type = $2
`

func (cl *ClientStore) GetVerificationEmail(db tutorme.DB, clientID string, emailType string) (string, error) {
	var email null.String
	err := db.QueryRowx(getVerificationEmail, clientID, emailType).Scan(&email)

	if err != nil {
		return "", err
	}

	if email.Valid == false {
		return "", errors.New("Could not find verification email")
	}

	return email.String, nil
}

const deleteVerificationEmail string = `
DELETE FROM email_verification
WHERE client_id = $1 AND email_type = $2
`

func (cl *ClientStore) DeleteVerificationEmail(db tutorme.DB, clientID string, emailType string) error {
	_, err := db.Queryx(deleteVerificationEmail, clientID, emailType)

	return err
}
