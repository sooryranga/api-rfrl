package store

import (
	"database/sql"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SessionStore holds all store related function for session
type SessionStore struct{}

// NeSessionStore creates new SessionStore
func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

type getSessionClientsResult struct {
	tutorme.Client
	SessionID int
}

const getSessionClients string = `
SELECT client.*, session_client.session_id FROM session_client
JOIN client ON session_client.client_id = client.id
WHERE session_client.session_id in (?)
	`

func getSessionWithClients(db tutorme.DB, rows *sqlx.Rows) (*[]tutorme.Session, error) {
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

	query, args, err := sqlx.In(getSessionClients, sessionIds)

	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	rows, err = db.Queryx(query, args...)
	var clients *[]tutorme.Client

	for rows.Next() {
		var result getSessionClientsResult

		err = rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		clients = &sessions[idToIndex[result.SessionID]].Clients
		*clients = append(*clients, result.Client)
	}

	return &sessions, nil
}

const getSessionByClientID string = `
SELECT .* FROM tutor_session
JOIN session_client ON session.id = session_clients.session_id
WHERE session_client.client_id = $1 AND session.state = $2
	`

func (ss *SessionStore) GetSessionByClientID(db tutorme.DB, clientID string, state string) (*[]tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByClientID, clientID, state)

	if err != nil {
		return nil, err
	}
	return getSessionWithClients(db, rows)
}

const getSessionByRoomID string = `
SELECT * FROM tutor_session
WHERE room_id = $1 AND state = $2
	`

func (ss *SessionStore) GetSessionByRoomID(db tutorme.DB, roomID string, state string) (*[]tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByRoomID, roomID, state)

	if err != nil {
		return nil, err
	}
	return getSessionWithClients(db, rows)
}

const checkSessionsIsForClient string = `
SELECT COUNT(*) from session_client
WHERE session_id in (?) and client_id = ?
	`

func (ss *SessionStore) CheckSessionsIsForClient(db tutorme.DB, clientID string, ids []int) (bool, error) {
	var m int
	query, args, err := sqlx.In(checkSessionsIsForClient, ids, clientID)
	if err != nil {
		return false, err
	}

	query = db.Rebind(query)

	row := db.QueryRowx(query, args...)
	err = row.Scan(&m)

	if err != nil {
		return false, err
	}

	return len(ids) == m, nil
}

const getSessionByID string = `
SELECT * FROM tutor_session
WHERE id = $1
	`

func (ss *SessionStore) GetSessionByID(db tutorme.DB, ID string) (*tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByID, ID)

	if err != nil {
		return nil, err
	}
	sessions, err := getSessionWithClients(db, rows)

	if err != nil {
		return nil, err
	}

	if len(*sessions) == 0 {
		return nil, errors.New("Session is not found")
	}

	return &(*sessions)[0], nil
}

const deleteSession string = `
DELETE FROM tutor_session
WHERE id = $1
	`

func (ss *SessionStore) DeleteSession(db tutorme.DB, ID int) error {
	_, err := db.Queryx(deleteSession, ID)

	return err
}

const insertSession string = `
INSERT INTO tutor_session (tutor_id, by, room_id)
VALUES ($1, $2, $3)
RETURNING *
	`

// CreateSession creates a new row for a document in the database
func (ss SessionStore) CreateSession(
	db tutorme.DB,
	session *tutorme.Session,
) (*tutorme.Session, error) {
	row := db.QueryRowx(
		insertSession,
		session.TutorID,
		session.UpdatedBy,
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

func (ss SessionStore) CreateClientSelectionOfEvent(
	db tutorme.DB,
	sessionID int,
	clientID string,
	canAttend bool,
) error {
	query := sq.Insert("client_selected_session").
		Columns("client_id", "session_id", "can_attend").
		Values(clientID, sessionID, canAttend)

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return err
}

const getSessionByIDForUpdate string = `
SELECT * FROM tutor_session
WHERE id = $1
FOR UPDATE OF session
	`

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
	eventID sql.NullInt64,
) (*tutorme.Session, error) {
	query := sq.Update("tutor_session").
		Set("by", by).
		Set("state", state)

	if eventID.Valid {
		query = query.Set("event_id", eventID)
	}

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
	query := sq.Insert("scheduled_event").
		Columns("start", "end", "title", "session_id")

	for i := 0; i < len(events); i++ {
		ev := events[i]
		query.Values(ev.Start, ev.End, ev.Title, ev.SessionID)
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

const getSessionEventFromSessionID string = `
SELECT scheduled_event.* FROM scheduled_event
JOIN tutor_session ON tutor_session.event_id = scheduled_event.id
WHERE tutor_session.id = ?
	`

func (ss SessionStore) GetSessionEventFromSessionID(
	db tutorme.DB,
	ID int,
) (*tutorme.Event, error) {
	row := db.QueryRowx(getSessionEventFromSessionID, ID)

	var event tutorme.Event

	err := row.StructScan(&event)

	return &event, err
}

const (
	getSessionEventFromClientIDs string = `
SELECT scheduled_event.* FROM tutor_session
JOIN session_client.session_id = tutor_session.id
JOIN scheduled_event ON tutor_session.event_id = scheduled_event.id
WHERE session_client.client_id IN (?) AND tutor_session.state = ?
	`

	getScheduledEventfromClientIDs string = `
SELECT scheduled_event.* FROM scheduled_event
JOIN client_event ON client_event.event_id = scheduled_event.id
WHERE client_event.client_id IN (?)
	`
)

func (ss SessionStore) GetScheduledEventsFromClientIDs(
	db tutorme.DB,
	clientIds []string,
	state string,
) (*[]tutorme.Event, error) {
	// Events related to session that is scheduled
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

	// Events related to user created event
	query, args, err = sqlx.In(getScheduledEventfromClientIDs, clientIds)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var e tutorme.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		*events = append(*events, e)
	}

	rows, err = db.Queryx(query, args...)

	return events, nil
}

const deleteSessionEvents string = `
DELETE FROM scheduled_event
WHERE id in (?) AND session_id = ?
	`

func (ss SessionStore) DeleteSessionEvents(db tutorme.DB, eventIDs []int, sessionId int) error {
	query, args, err := sqlx.In(deleteSessionEvents, eventIDs, sessionId)

	if err != nil {
		return err
	}
	query = db.Rebind(query)

	_, err = db.Queryx(query, args...)

	return err
}

const getSessionEventByID string = `
SELECT scheduled_event.* from scheduled_event
JOIN tutor_session ON tutor_session.event_id = scheduled_event.id
WHERE scheduled_event.id = $1 AND tutor_session.id = $2
	`

func (ss SessionStore) GetSessionEventByID(db tutorme.DB, sessionID int, ID int) (*tutorme.Event, error) {
	row := db.QueryRowx(getSessionEventByID, ID, sessionID)

	var event tutorme.Event

	err := row.StructScan(&event)

	if err != nil {
		return nil, err
	}

	return &event, err
}

func (ss SessionStore) CheckOverlapingEvents(db tutorme.DB, ID int, events *[]tutorme.Event) (bool, error) {
	query := sq.Select("*").
		Prefix("SELECT EXISTS(").
		From("scheduled_events")

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
