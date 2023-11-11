package store

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

// SessionStore holds all store related function for session
type SessionStore struct{}

// NeSessionStore creates new SessionStore
func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

type getSessionClientsResult struct {
	tutorme.Client
	CanAttend null.Bool `db:"can_attend"`
	SessionID int       `db:"session_id"`
}

const getSessionClients string = `
SELECT 
	client.*, 
	session_client.session_id as "session_id",
	session_client.can_attend as "can_attend" 
FROM session_client
JOIN client ON session_client.client_id = client.id
WHERE session_client.session_id in (?)
	`

func getSessionWithClients(db tutorme.DB, rows *sqlx.Rows, clientID string) (*[]tutorme.Session, error) {
	idToIndex := make(map[int]int)
	sessions := make([]tutorme.Session, 0)
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
		return &sessions, nil
	}

	query, args, err := sqlx.In(getSessionClients, sessionIds)

	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	rows, err = db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var result getSessionClientsResult

		err = rows.StructScan(&result)
		if err != nil {
			return nil, err
		}

		index := idToIndex[result.SessionID]
		session := &sessions[index]

		if session.TutorID == result.Client.ID {
			session.Tutor = result.Client
		}

		if result.Client.ID == clientID {
			session.CanAttend = result.CanAttend
		}

		session.Clients = append(session.Clients, result.Client)
	}

	return &sessions, nil
}

func (ss *SessionStore) GetSessionByClientID(db tutorme.DB, clientID string, state string) (*[]tutorme.Session, error) {
	query := sq.
		Select(`tutor_session.*`).
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

	return getSessionWithClients(db, rows, clientID)
}

const getSessionByRoomID string = `
SELECT * FROM tutor_session
WHERE room_id = $1 AND state = $2
	`

func (ss *SessionStore) GetSessionByRoomID(db tutorme.DB, clientID string, roomID string, state string) (*[]tutorme.Session, error) {
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

	return getSessionWithClients(db, rows, clientID)
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

func (ss *SessionStore) GetSessionByID(db tutorme.DB, clientID string, ID int) (*tutorme.Session, error) {
	rows, err := db.Queryx(getSessionByID, ID)

	if err != nil {
		return nil, err
	}

	sessions, err := getSessionWithClients(db, rows, clientID)

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

const getSessionsEvent string = `
SELECT scheduled_event.*, tutor_session.id as "session_id" FROM tutor_session
INNER JOIN scheduled_event ON tutor_session.event_id = scheduled_event.id
WHERE tutor_session.id in (?)
`

type getSessionsEventResult struct {
	tutorme.Event
	SessionID int `db:"session_id"`
}

func (ss SessionStore) GetSessionsEvent(db tutorme.DB, sessionIDs []int) (map[int]*tutorme.Event, error) {
	sessionIDToEventMap := make(map[int]*tutorme.Event)
	query, args, err := sqlx.In(getSessionsEvent, sessionIDs)

	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var result getSessionsEventResult

		err := rows.StructScan(&result)

		if err != nil {
			return nil, err
		}
		sessionIDToEventMap[result.SessionID] = &result.Event
	}

	return sessionIDToEventMap, nil
}

const getSessionByIDForUpdate string = `
SELECT 
	tutor_session.*, 
	session_client.can_attend as "can_attend" 
FROM tutor_session
INNER JOIN session_client ON session_client.session_id = tutor_session.id
WHERE tutor_session.id = $1
FOR UPDATE OF tutor_session
	`

func (ss SessionStore) GetSessionByIDForUpdate(db tutorme.DB, clientID string, ID int) (*tutorme.Session, error) {
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
	eventID null.Int,
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

	insertedEvents := make([]tutorme.Event, 0)

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

const checkAllClientSessionHasRespondedQuery string = `
SELECT EXISTS(
	SELECT 1 FROM session_client
	WHERE session_client.session_id = $1 AND 
	session_client.can_attend = NULL
)
`

func (ss SessionStore) CheckAllClientSessionHasResponded(db tutorme.DB, id int) (bool, error) {
	var notAllClientsResponded null.Bool
	err := db.QueryRowx(checkAllClientSessionHasRespondedQuery, id).Scan(&notAllClientsResponded)

	if err != nil {
		return false, err
	}

	if !notAllClientsResponded.Valid {
		return false, errors.New("Unexpected invalid bool returned from database")
	}

	return !notAllClientsResponded.Bool, nil
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
