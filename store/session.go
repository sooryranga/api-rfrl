package store

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SessionStore holds all store relateed function for session
type SessionStore struct{}

// NeSessionStore creates new SessionStore
func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

const (
	getSessionByClientID string = `
SELECT sesssion.* FROM session
JOIN session_client ON session.id = session_clients.session_id
WHERE session_client.client_id = $1 AND session.state = $2
	`
	getSessionByRoomID string = `
SELECT * FROM session
WHERE room_id = $1 AND state = $2
	`
	getSessionByID string = `
SELECT * FROM session
WHERE id = $1
	`
	getSessionByIDForUpdate string = `
SELECT * FROM session
WHERE id = $1
FOR UPDATE OF session
	`
	checkSessionsIsForClient string = `
SELECT COUNT(*) from session_client
WHERE session_id in (?) and client_id = ?
	`
	insertSession string = `
INSERT INTO session (tutor_id, by, room_id)
VALUES ($1, $2, $3)
RETURNING *
	`
	deleteSession string = `
DELETE FROM session
WHERE id = $1
	`
	getSessionCients string = `
SELECT client.*, session_client.session_id FROM session_client
JOIN client ON session_client.client_id = client.id
WHERE session_client.session_id in (?)
	`
	getSessionEventFromSessionID string = `
SELECT * FROM session_event
WHERE session_id = $1
	`
	getSessionEventFromClientIDs string = `
SELECT session_event FROM session_event 
JOIN session ON session.id = session_event.session_id
JOIN session_client.session_id = session_event.id
WHERE session_client.client_id in (?) AND session.state = ?
	`
	deleteSessionEvents string = `
DELETE FROM session_event
WHERE id in (?)
	`
)

type getSessionCientsResult struct {
	tutorme.Client
	SessionId int
}

func getSessionWithClientsById(db tutorme.DB, rows *sqlx.Rows) (*[]tutorme.Session, error) {
	var idToIndex map[int]int
	var sessions []tutorme.Session
	var sessionIds []int
	i := 0
	for rows.Next() {
		var session tutorme.Session
		err := rows.StructScan(&session)

		if err != nil {
			return nil, err
		}
		idToIndex[session.ID] = i
		sessions = append(sessions, session)
		sessionIds = append(sessionIds, session.ID)
		i++
	}

	query, args, err := sqlx.In(getSessionCients, sessionIds)

	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	rows, err = db.Queryx(query, args...)
	var clients *[]tutorme.Client

	for rows.Next() {
		var result getSessionCientsResult

		err = rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		clients = &sessions[idToIndex[result.SessionId]].Clients
		*clients = append(*clients, result.Client)
	}

	return &sessions, nil
}

func (ss *SessionStore) GetSessionByClientID(db tutorme.DB, clientID string, state string) (*[]tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByClientID, clientID, state)

	if err != nil {
		return nil, err
	}
	return getSessionWithClientsById(db, rows)
}

func (ss *SessionStore) GetSessionByRoomID(db tutorme.DB, roomID string, state string) (*[]tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByRoomID, roomID, state)

	if err != nil {
		return nil, err
	}
	return getSessionWithClientsById(db, rows)
}

func (ss *SessionStore) CheckSessionsIsForClient(db tutorme.DB, clientID string, ids []int) (bool, error) {
	var m int
	query, args, err := sqlx.In(checkSessionsIsForClient, ids, clientID)
	if err != nil {
		return false, err
	}

	query = db.Rebind(query)

	row := db.QueryRowx(query, args...)
	err := row.Scan(&m)

	if err != nil {
		return false, err
	}

	return len(ids) == m, nil
}

func (ss *SessionStore) GetSessionByID(db tutorme.DB, ID string) (*tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByID, ID)

	if err != nil {
		return nil, err
	}
	sessions, err := getSessionWithClientsById(db, rows)

	if err != nil {
		return nil, err
	}

	if len(*sessions) == 0 {
		return nil, errors.New("Session is not found")
	}

	return &(*sessions)[0], nil
}

func (ss *SessionStore) DeleteSession(db tutorme.DB, ID int) error {
	_, err := db.Queryx(deleteSession, ID)

	return err
}

// CreateSession creates a new row for a document in the database
func (ss SessionStore) CreateSession(
	db tutorme.DB,
	session *tutorme.Session,
) (*tutorme.Session, error) {
	row := db.QueryRowx(
		insertSession,
		session.TutorID,
		session.By,
		session.RoomID,
	)

	var m tutorme.Session

	err := row.StructScan(&m)

	return &m, errors.Wrap(err, "CreateSession")
}

func (ss SessionStore) CreateSessionClients(
	db tutorme.DB,
	sessionID int,
	clientIDs []string,
) (*[]tutorme.Client, error) {
	query := sq.Insert("session_client").Columns("session_id, client_id")

	for i := 0; i < len(clientIDs); i++ {
		query.Values(sessionID, clientIDs[i])
	}
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return getClientFromIDs(db, clientIDs)
}

func (ss SessionStore) GetSessionByIDForUpdate(db tutorme.DB, id int) (tutorme.Session, error) {
	row := db.QueryRowx(getSessionByIDForUpdate, id)

	var m tutorme.Session

	err := row.StructScan(&m)
	return m, err
}

func (ss SessionStore) UpdateSession(
	db tutorme.DB,
	id int,
	by string,
	state string,
) (*tutorme.Session, error) {
	query := sq.Update("session").
		Set("by", by).
		Set("state", state)

	sql, args, err := query.
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(sql, args...)

	var m tutorme.Session

	err = row.StructScan(&m)

	return &m, err
}

func (ss SessionStore) CreateSessionEvents(
	db tutorme.DB,
	events []tutorme.Event,
) (*[]tutorme.Event, error) {
	query := sq.Insert("session_event").
		Columns("start", "end", "title", "session_id")

	for i := 0; i < len(events); i++ {
		ev := events[i]
		query.Values(ev.Start, ev.End, ev.Title, ev.SessionId)
	}
	sql, args, err := query.
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(
		sql,
		args...,
	)

	var insertedEvents *[]tutorme.Event

	for rows.Next() {
		var e tutorme.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		*insertedEvents = append(*insertedEvents, e)
	}

	return insertedEvents, err
}

func (ss SessionStore) GetSessionEventFromSessionID(
	db tutorme.DB,
	ID int,
) (*[]tutorme.Event, error) {
	rows, err := db.Queryx(getSessionEventFromSessionID, ID)

	if err != nil {
		return nil, err
	}

	var events *[]tutorme.Event

	for rows.Next() {
		var e tutorme.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		*events = append(*events, e)
	}
	return events, nil
}

func (ss SessionStore) GetSessionEventFromClientIDs(
	db tutorme.DB,
	clientIds []string,
	state string,
) (*[]tutorme.Event, error) {
	query, args, err := sqlx.In(getSessionEventFromClientIDs, clientIds, state)

	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	var events *[]tutorme.Event

	for rows.Next() {
		var e tutorme.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		*events = append(*events, e)
	}
	return events, nil
}

func (ss SessionStore) DeleteSessionEvents(db tutorme.DB, IDs []int) error {
	query, args, err := sqlx.In(deleteSessionEvents, IDs)

	if err != nil {
		return err
	}
	query = db.Rebind(query)

	_, err = db.Queryx(query, args...)

	return err
}

func (ss SessionStore) CheckOverlapingEvents(db tutorme.DB, ID int, events *[]tutorme.Event) (bool, error) {
	query := sq.Select("*").
		Prefix("SELECT EXISTS(").
		From("session_events")

	for i := 0; i < len(*events); i++ {
		event := (*events)[i]
		query = query.Where(
			sq.Or{
				sq.And{
					sq.LtOrEq{"start": event.Start},
					sq.Gt{"end": event.Start},
					sq.Eq{"session_id": ID},
				},
				sq.And{
					sq.Lt{"start": event.End},
					sq.GtOrEq{"end": event.End},
					sq.Eq{"session_id": ID},
				},
				sq.And{
					sq.GtOrEq{"start": event.Start},
					sq.LtOrEq{"end": event.End},
					sq.Eq{"session_id": ID},
				},
			},
		)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return true, err
	}

	var m bool

	row := db.QueryRowx(sql, args...)
	err = row.Scan(&m)

	return m, err
}
