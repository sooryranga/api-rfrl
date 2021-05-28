package store

import (
	"database/sql"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
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
	query := sq.Select("client.*").From("client")

	if options.IsTutor.Valid {
		query = query.Where(sq.Eq{"is_tutor": options.IsTutor.Bool})
	}

	if len(options.CompanyIds) > 0 {
		query = query.Where(sq.Eq{"company_id": options.CompanyIds})
	}

	if options.WantingReferralCompanyId.Valid {
		query = query.Join("client_wanting_company_referral ON client.id = client_wanting_company_referral.client_id").
			Where(sq.Eq{"client_wanting_company_referral.company_id": options.WantingReferralCompanyId.Int64})
	}

	sql, args, err := query.
		OrderBy("id DESC").
		PlaceholderFormat(sq.Dollar).ToSql()

	log.Error(sql)

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
	if client.VerifiedEmail.Valid {
		query = query.Set("verified_email", client.VerifiedEmail)
	}
	if client.IsLookingForReferral.Valid {
		query = query.Set("is_looking_for_referral", client.IsLookingForReferral)
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
ON CONFLICT ON CONSTRAINT email_verification_client_id_email_type_key
DO UPDATE
SET email = $2, pass_code = $4
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
	SELECT id, pass_code
	FROM email_verification
	WHERE client_id = $1
	AND email = $2
	AND email_type = $3
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
	var expectedPasscode string
	err := db.QueryRowx(verifyEmail, clientID, email, emailType).Scan(&id, &expectedPasscode)

	if err != nil {
		return err
	}

	if !id.Valid {
		return errors.New("Could not find email to verify")
	}

	if expectedPasscode != passCode {
		return errors.New("Passcode is not correct")
	}

	_, err = db.Queryx(deleteVerifyEmailFromId, id)

	return err
}

const getVerificationEmailQuery string = `
SELECT email FROM email_verification
WHERE client_id = $1 AND email_type = $2
`

func (cl *ClientStore) GetVerificationEmail(db tutorme.DB, clientID string, emailType string) (string, error) {
	var email null.String
	err := db.QueryRowx(getVerificationEmailQuery, clientID, emailType).Scan(&email)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Could not find verification email")
		}
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

func (cl ClientStore) GetRelatedEventsByClientIDs(
	db tutorme.DB,
	clientIDs []string,
	start null.Time,
	end null.Time,
	state null.String,
) (*[]tutorme.Event, error) {
	sessionQuery := getSessionEventsRelatedToClientsQuery(clientIDs)

	if start.Valid {
		sessionQuery = sessionQuery.Where(sq.GtOrEq{"scheduled_event.start_time": start})
	}

	if end.Valid {
		sessionQuery = sessionQuery.Where(sq.LtOrEq{"scheduled_event.end_time": end})
	}

	if state.Valid {
		sessionQuery = sessionQuery.Where(sq.Eq{"tutor_session.state": state})
	}

	events := make([]tutorme.Event, 0)
	sql, args, err := sessionQuery.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return &events, err
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return &events, err
	}

	for rows.Next() {
		var event tutorme.Event

		err = rows.StructScan(&event)

		if err != nil {
			return &events, err
		}
		events = append(events, event)
	}

	clientQuery := getEventsRelatedToClientsQuery(clientIDs)

	if start.Valid {
		clientQuery = clientQuery.Where(sq.GtOrEq{"start_time": start})
	}

	if end.Valid {
		clientQuery = clientQuery.Where(sq.LtOrEq{"end_time": end})
	}

	sql, args, err = clientQuery.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return &events, err
	}

	rows, err = db.Queryx(sql, args...)

	if err != nil {
		return &events, err
	}

	for rows.Next() {
		var event tutorme.Event

		err = rows.StructScan(&event)

		if err != nil {
			return &events, err
		}
		events = append(events, event)
	}

	return &events, nil
}

func getSessionEventsRelatedToClientsQuery(clientIDs []string) sq.SelectBuilder {
	return sq.Select("scheduled_event.*").
		From("scheduled_event").
		Join("tutor_session ON tutor_session.event_id = scheduled_event.id").
		Join("session_client ON session_client.session_id = tutor_session.id").
		Where(sq.Eq{"session_client.client_id": clientIDs})
}

func getEventsRelatedToClientsQuery(clientIDs []string) sq.SelectBuilder {
	return sq.Select("*").
		From("scheduled_event").
		Join("client_event ON client_event.event_id = scheduled_event.id").
		Where(sq.Eq{"client_event.client_id": clientIDs})
}

func checkOverlapingSessionEvents(db tutorme.DB, clientIDs []string, events *[]tutorme.Event) (bool, error) {
	query := getSessionEventsRelatedToClientsQuery(clientIDs).
		Where(sq.Eq{"tutor_session.state": tutorme.SCHEDULED}).
		Prefix("SELECT EXISTS(")

	query = filterInclusiveDateRange(query, events)

	sql, args, err := query.Suffix(")").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return true, err
	}

	m := false

	row := db.QueryRowx(sql, args...)
	err = row.Scan(&m)

	return m, err
}

func checkOverlapingClientEvents(db tutorme.DB, clientIDs []string, events *[]tutorme.Event) (bool, error) {
	query := getEventsRelatedToClientsQuery(clientIDs).
		Prefix("SELECT EXISTS(")

	query = filterInclusiveDateRange(query, events)

	sql, args, err := query.Suffix(")").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return true, err
	}

	m := false

	row := db.QueryRowx(sql, args...)
	err = row.Scan(&m)

	return m, err
}

func (cl ClientStore) CheckOverlapingEventsByClientIDs(db tutorme.DB, clientIds []string, events *[]tutorme.Event) (bool, error) {
	clientOverlap, err := checkOverlapingClientEvents(db, clientIds, events)

	if err != nil {
		return true, err
	}
	sessionOverlap, err := checkOverlapingSessionEvents(db, clientIds, events)

	if err != nil {
		return true, err
	}

	log.Errorj(log.JSON{
		"clientOverlap":  clientOverlap,
		"sessionOverlap": sessionOverlap,
	})

	return clientOverlap || sessionOverlap, nil
}

func (cl ClientStore) CreateOrUpdateClientEducation(db tutorme.DB, clientID string, education tutorme.Education) error {
	query := sq.Update("client")

	if education.Institution.Valid {
		query = query.Set("institution", education.Institution)
	}
	if education.Degree.Valid {
		query = query.Set("degree", education.Degree)
	}
	if education.FieldOfStudy.Valid {
		query = query.Set("field_of_study", education.FieldOfStudy)
	}
	if education.StartYear.Valid {
		query = query.Set("start_year", education.StartYear)
	}
	if education.EndYear.Valid {
		query = query.Set("end_year", education.EndYear)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": clientID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return err
}

const deleteClientWantingCompanyReferralsSQuery string = `
DELETE FROM client_wanting_company_referral WHERE client_id = $1
`

func (cl ClientStore) CreateClientWantingCompanyReferrals(db tutorme.DB, clientID string, companyIDs []int) error {
	_, err := db.Queryx(deleteClientWantingCompanyReferralsSQuery, clientID)

	if err != nil {
		return err
	}

	if len(companyIDs) == 0 {
		return nil
	}

	query := sq.Insert("client_wanting_company_referral").
		Columns("client_id", "company_id")

	for i := 0; i < len(companyIDs); i++ {
		query = query.Values(clientID, companyIDs[i])
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	_, err = db.Queryx(sql, args...)

	return err
}

const getClientWantingCompanyReferralsQuery string = `
SELECT company_id FROM client_wanting_company_referral WHERE client_id = $1
`

func (cl ClientStore) GetClientWantingCompanyReferrals(db tutorme.DB, clientID string) ([]int, error) {
	companyIDs := make([]int, 0)

	rows, err := db.Queryx(getClientWantingCompanyReferralsQuery, clientID)

	if err != nil {
		return companyIDs, err
	}

	for rows.Next() {
		var companyID int
		err = rows.Scan(&companyID)
		if err != nil {
			return companyIDs, err
		}
		companyIDs = append(companyIDs, companyID)
	}

	return companyIDs, nil
}
