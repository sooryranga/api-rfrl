package store

import (
	"database/sql"

	"github.com/Arun4rangan/api-tutorme/tutorme"
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

func (cs *ConferenceStore) GetOrCreateConference(db tutorme.DB, sessionID int) (*tutorme.Conference, error) {
	var conference tutorme.Conference
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

func (cs *ConferenceStore) CreateNewCode(db tutorme.DB, sessionID int, rawCode string) (*tutorme.Code, error) {
	var code tutorme.Code
	row := db.QueryRowx(createCodeQuery, rawCode)

	err := row.StructScan(&code)

	if err != nil {
		return nil, err
	}

	_, err = db.Queryx(updateCodeConferenceQuery, sessionID, code.ID)

	return &code, err
}
