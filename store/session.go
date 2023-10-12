package store

import (
	"database/sql"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
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
	SessionID int `db:"session_id"`
}

const getSessionClients string = `
SELECT client.*, session_client.session_id as "session_id" FROM session_client
JOIN client ON session_client.client_id = client.id
WHERE session_client.session_id in (?)
	`

func getSessionWithClients(db tutorme.DB, rows *sqlx.Rows) (*[]tutorme.Session, error) {
	idToIndex := make(map[int]int)
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

	if i == 0 {
		return nil, nil
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

func (ss *SessionStore) GetSessionByClientID(db tutorme.DB, clientID string, state string) (*[]tutorme.Session, error) {
	query := sq.
		Select("tutor_session.*").
		From("tutor_session").
		Join("session_client ON tutor_session.id = session_client.session_id").
		Where(sq.Eq{"session_client.client_id": clientID})

	if state != "" {
		query = query.Where(sq.Eq{"state": state})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(sql, args...)

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
	query := sq.Select("*").From("tutor_session").Where(sq.Eq{"room_id": roomID})

	if state != "" {
		query = query.Where(sq.Eq{"state": state})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(sql, args...)

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
	log.Errorj(log.JSON{"count": m, "err": err, "session_id": ids, "clientID": clientID})
	if err != nil {
		return false, err
	}

	return len(ids) == m, nil
}

const getSessionByID string = `
SELECT * FROM tutor_session
WHERE id = $1
	`

func (ss *SessionStore) GetSessionByID(db tutorme.DB, ID int) (*tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByID, ID)

	if err != nil {
		return nil, err
	}

	sessions, err := getSessionWithClients(db, rows)

	if err != nil {
		return nil, err
	}

	if sessions == nil || len(*sessions) == 0 {
		return nil, errors.New("Session is not found")
	}

	return &((*sessions)[0]), nil
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
INSERT INTO tutor_session (tutor_id, updated_by, room_id, state)
VALUES ($1, $2, $3, $4)
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
		session.State,
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
		query = query.Values(sessionID, clientIDs[i])
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
	query := sq.Update("session_client").
		Set("can_attend", canAttend).
		Where(sq.Eq{"session_id": sessionID}).
		Where(sq.Eq{"client_id": clientID})

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
FOR UPDATE OF tutor_session
	`

func (ss SessionStore) GetSessionByIDForUpdate(db tutorme.DB, ID int) (*tutorme.Session, error) {
	row := db.QueryRowx(getSessionByIDForUpdate, ID)

	var m tutorme.Session

	err := row.StructScan(&m)
	return &m, err
}

func (ss SessionStore) UpdateSession(
	db tutorme.DB,
	id int,
	by string,
	state string,
	eventID sql.NullInt64,
) (*tutorme.Session, error) {
	query := sq.Update("tutor_session").Set("updated_by", by)

	if state != "" {
		query = query.Set("state", state)
	}

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
		Columns("start_time", "end_time", "title")

	for i := 0; i < len(events); i++ {
		ev := events[i]
		query = query.Values(ev.StartTime, ev.EndTime, ev.Title)
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

	if err != nil {
		return nil, err
	}

	var insertedEvents []tutorme.Event

	for rows.Next() {
		var e tutorme.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		insertedEvents = append(insertedEvents, e)
	}

	return &insertedEvents, err
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

const getScheduledEventfromClientIDs string = `
SELECT scheduled_event.* FROM scheduled_event
JOIN client_event ON client_event.event_id = scheduled_event.id
WHERE client_event.client_id IN (?)
	`

func (ss SessionStore) GetScheduledEventsFromClientIDs(
	db tutorme.DB,
	clientIds []string,
	state string,
) (*[]tutorme.Event, error) {
	// Events related to session that is scheduled
	sessionQuery := sq.Select("scheduled_event.*").
		From("tutor_session").
		Join("session_client ON session_client.session_id = tutor_session.id").
		Join("scheduled_event ON tutor_session.event_id = scheduled_event.id").
		Where(sq.Eq{"session_client.client_id": clientIds}).
		Where(sq.Eq{"tutor_session.state": tutorme.SCHEDULED})

	if state != "" {
		sessionQuery = sessionQuery.Where(sq.Eq{"tutor_session.state": state})
	}
	sql, args, err := sessionQuery.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(sql, args...)

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
	clientQuery, args, err := sqlx.In(getScheduledEventfromClientIDs, clientIds)

	if err != nil {
		return nil, err
	}

	clientQuery = db.Rebind(clientQuery)

	rows, err = db.Queryx(clientQuery, args...)

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

	return events, nil
}

const deleteSessionEvents string = `
DELETE FROM scheduled_event
WHERE id in (?)
	`

func (ss SessionStore) DeleteSessionEvents(db tutorme.DB, eventIDs []int) error {
	query, args, err := sqlx.In(deleteSessionEvents, eventIDs)

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

func filterInclusiveDateRange(query sq.SelectBuilder, events *[]tutorme.Event) sq.SelectBuilder {
	for i := 0; i < len(*events); i++ {
		event := (*events)[i]
		query = query.Where(
			sq.Or{
				sq.And{
					sq.LtOrEq{"start_time": event.StartTime},
					sq.Gt{"end_time": event.StartTime},
				},
				sq.And{
					sq.Lt{"start_time": event.EndTime},
					sq.GtOrEq{"end_time": event.EndTime},
				},
				sq.And{
					sq.GtOrEq{"start_time": event.StartTime},
					sq.LtOrEq{"end_time": event.EndTime},
				},
			},
		)
	}
	return query
}

func checkOverlapingSessionEvents(db tutorme.DB, clientIDs []string, events *[]tutorme.Event) (bool, error) {
	query := sq.Select("*").
		Prefix("SELECT EXISTS(").
		From("scheduled_event").
		Join("tutor_session ON tutor_session.event_id = scheduled_event.id").
		Join("session_client ON session_client.session_id = tutor_session.id").
		Where(sq.Eq{"session_client.client_id": clientIDs}).
		Where(sq.Eq{"tutor_session.state": tutorme.SCHEDULED})

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
	query := sq.Select("*").
		Prefix("SELECT EXISTS(").
		From("scheduled_event").
		Join("client_event ON client_event.event_id = scheduled_event.id").
		Where(sq.Eq{"client_event.client_id": clientIDs})

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

func (ss SessionStore) CheckOverlapingEvents(db tutorme.DB, clientIds []string, events *[]tutorme.Event) (bool, error) {
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

const checkClientsAttendedTutorSession string = `
SELECT count(client_id) FROM session_client
JOIN tutor_session ON tutor_session.id = session_client.session_id
WHERE client_id IN (?) 
AND can_attend = TRUE 
AND tutor_session.state = 'paid'
AND tutor_session.tutor_id = ?
GROUP BY session_id
	`

func (ss SessionStore) CheckClientsAttendedTutorSession(db tutorme.DB, tutorID string, clientIDs []string) (bool, error) {
	sql, args, err := sqlx.In(checkClientsAttendedTutorSession, clientIDs, tutorID)

	if err != nil {
		return false, err
	}

	sql = db.Rebind(sql)

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return false, err
	}

	for rows.Next() {
		var b int
		err := rows.Scan(&b)

		if err != nil {
			return false, err
		}

		if b == len(clientIDs) {
			return true, nil
		}
	}

	return false, nil
}
