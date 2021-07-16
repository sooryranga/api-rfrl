package store

import (
	"database/sql"

	"github.com/Arun4rangan/api-rfrl/rfrl"
	sq "github.com/Masterminds/squirrel"
)

type ConferenceStore struct{}

func NewConferenceStore() *ConferenceStore {
	return &ConferenceStore{}
}

const selectConferenceQuery = `
	SELECT * FROM session_conference
	WHERE session_id = $1 
`

const createConferenceQuery = `
INSERT INTO session_conference(session_id, latest_code)
VALUES ($1, null)
RETURNING *
`

func (cs *ConferenceStore) GetOrCreateConference(db rfrl.DB, sessionID int) (*rfrl.Conference, error) {
	var conference rfrl.Conference
	err := db.QueryRowx(selectConferenceQuery, sessionID).StructScan(&conference)

	if err != nil && err == sql.ErrNoRows {
		err = db.QueryRowx(createConferenceQuery, sessionID).StructScan(&conference)
	}

	return &conference, err
}

const createCodeQuery = `
INSERT INTO conference_code(code)
VALUES ($1)
RETURNING *
`

const updateCodeConferenceQuery = `
UPDATE session_conference
SET latest_code = $1
WHERE session_id = $2
`

func (cs *ConferenceStore) CreateNewCode(db rfrl.DB, sessionID int, rawCode string) (*rfrl.Code, error) {
	var code rfrl.Code
	row := db.QueryRowx(createCodeQuery, rawCode)

	err := row.StructScan(&code)

	if err != nil {
		return nil, err
	}

	_, err = db.Queryx(updateCodeConferenceQuery, sessionID, code.ID)

	return &code, err
}

func (cs *ConferenceStore) UpdateCode(db rfrl.DB, id int, code rfrl.Code) (*rfrl.Code, error) {
	query := sq.Update("conference_code")

	if code.Code.Valid {
		query = query.Set("code", code.Code.String)
	}

	if code.Result.Valid {
		query = query.Set("result", code.Result.String)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": id}).
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

	var c rfrl.Code

	err = row.StructScan(&c)

	return &c, err
}
