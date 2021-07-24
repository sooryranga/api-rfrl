package store

import (
	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
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
	rfrl.Client
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

func getSessionWithClients(db rfrl.DB, rows *sqlx.Rows, clientID string) (*[]rfrl.Session, error) {
	idToIndex := make(map[int]int)
	sessions := make([]rfrl.Session, 0)
	sessionIds := make([]int, 0)
	i := 0
	for rows.Next() {
		var session rfrl.Session
		err := rows.StructScan(&session)

		if err != nil {
			return nil, errors.Wrap(err, "getSessionWithClients")
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
		return nil, errors.Wrap(err, "getSessionWithClients")
	}

	query = db.Rebind(query)
	rows, err = db.Queryx(query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "getSessionWithClients")
	}

	for rows.Next() {
		var result getSessionClientsResult

		err = rows.StructScan(&result)
		if err != nil {
			return nil, errors.Wrap(err, "getSessionWithClients")
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

func (ss *SessionStore) GetSessionByClientID(db rfrl.DB, clientID string, state string) (*[]rfrl.Session, error) {
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
		return nil, errors.Wrap(err, "GetSessionByClientID")
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByClientID")
	}

	return getSessionWithClients(db, rows, clientID)
}

const getSessionByRoomID string = `
SELECT * FROM tutor_session
WHERE room_id = $1 AND state = $2
	`

func (ss *SessionStore) GetSessionByRoomID(db rfrl.DB, clientID string, roomID string, state string) (*[]rfrl.Session, error) {
	query := sq.Select("*").From("tutor_session").Where(sq.Eq{"room_id": roomID})

	if state != "" {
		query = query.Where(sq.Eq{"state": state})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByRoomID")
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByRoomID")
	}

	return getSessionWithClients(db, rows, clientID)
}

const checkSessionsIsForClient string = `
SELECT COUNT(*) from session_client
WHERE session_id in (?) and client_id = ?
	`

func (ss *SessionStore) CheckSessionsIsForClient(db rfrl.DB, clientID string, ids []int) (bool, error) {
	var m int
	query, args, err := sqlx.In(checkSessionsIsForClient, ids, clientID)
	if err != nil {
		return false, errors.Wrap(err, "CheckSessionsIsForClient")
	}

	query = db.Rebind(query)

	row := db.QueryRowx(query, args...)
	err = row.Scan(&m)

	return len(ids) == m, errors.Wrap(err, "CheckSessionsIsForClient")
}

const getSessionByID string = `
SELECT * FROM tutor_session
WHERE id = $1
	`

func (ss *SessionStore) GetSessionByID(db rfrl.DB, clientID string, ID int) (*rfrl.Session, error) {
	rows, err := db.Queryx(getSessionByID, ID)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByID")
	}

	sessions, err := getSessionWithClients(db, rows, clientID)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByID")
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

func (ss *SessionStore) DeleteSession(db rfrl.DB, ID int) error {
	_, err := db.Queryx(deleteSession, ID)

	return errors.Wrap(err, "DeleteSession")
}

const insertSession string = `
INSERT INTO tutor_session (tutor_id, updated_by, room_id, state)
VALUES ($1, $2, $3, $4)
RETURNING *
	`

// CreateSession creates a new row for a document in the database
func (ss SessionStore) CreateSession(
	db rfrl.DB,
	session *rfrl.Session,
) (*rfrl.Session, error) {
	row := db.QueryRowx(
		insertSession,
		session.TutorID,
		session.UpdatedBy,
		session.RoomID,
		session.State,
	)

	var m rfrl.Session

	err := row.StructScan(&m)

	return &m, errors.Wrap(err, "CreateSession")
}

func (ss SessionStore) CreateSessionClients(
	db rfrl.DB,
	sessionID int,
	clientIDs []string,
) (*[]rfrl.Client, error) {
	query := sq.Insert("session_client").Columns("session_id, client_id")

	for i := 0; i < len(clientIDs); i++ {
		query = query.Values(sessionID, clientIDs[i])
	}
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "CreateSessionClients")
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	if err != nil {
		return nil, errors.Wrap(err, "CreateSessionClients")
	}

	return getClientFromIDs(db, clientIDs)
}

func (ss SessionStore) CreateClientSelectionOfEvent(
	db rfrl.DB,
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
		return errors.Wrap(err, "CreateClientSelectionOfEvent")
	}

	_, err = db.Queryx(
		sql,
		args...,
	)

	return errors.Wrap(err, "CreateClientSelectionOfEvent")
}

const getSessionsEvent string = `
SELECT scheduled_event.*, tutor_session.id as "session_id" FROM tutor_session
INNER JOIN scheduled_event ON tutor_session.event_id = scheduled_event.id
WHERE tutor_session.id in (?)
`

type getSessionsEventResult struct {
	rfrl.Event
	SessionID int `db:"session_id"`
}

func (ss SessionStore) GetSessionsEvent(db rfrl.DB, sessionIDs []int) (map[int]*rfrl.Event, error) {
	sessionIDToEventMap := make(map[int]*rfrl.Event)
	query, args, err := sqlx.In(getSessionsEvent, sessionIDs)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionsEvent")
	}
	query = db.Rebind(query)

	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionsEvent")
	}

	for rows.Next() {
		var result getSessionsEventResult

		err := rows.StructScan(&result)

		if err != nil {
			return nil, errors.Wrap(err, "GetSessionsEvent")
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

func (ss SessionStore) GetSessionByIDForUpdate(db rfrl.DB, clientID string, ID int) (*rfrl.Session, error) {
	row := db.QueryRowx(getSessionByIDForUpdate, ID)

	var m rfrl.Session

	err := row.StructScan(&m)
	return &m, errors.Wrap(err, "GetSessionByIDForUpdate")
}

func (ss SessionStore) UpdateSession(
	db rfrl.DB,
	id int,
	by string,
	state string,
	eventID null.Int,
	conferenceID null.String,
) (*rfrl.Session, error) {
	query := sq.Update("tutor_session").Set("updated_by", by)

	if state != "" {
		query = query.Set("state", state)
	}

	if eventID.Valid {
		query = query.Set("event_id", eventID)
	}

	if conferenceID.Valid {
		query = query.Set("conference_id", conferenceID)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "UpdateSession")
	}

	row := db.QueryRowx(sql, args...)

	var m rfrl.Session

	err = row.StructScan(&m)

	return &m, errors.Wrap(err, "UpdateSession")
}

func (ss SessionStore) CreateSessionEvents(
	db rfrl.DB,
	events []rfrl.Event,
) (*[]rfrl.Event, error) {
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
		return nil, errors.Wrap(err, "CreateSessionEvents")
	}

	rows, err := db.Queryx(
		sql,
		args...,
	)

	if err != nil {
		return nil, errors.Wrap(err, "CreateSessionEvents")
	}

	insertedEvents := make([]rfrl.Event, 0)

	for rows.Next() {
		var e rfrl.Event

		err = rows.StructScan(&e)
		if err != nil {
			return nil, errors.Wrap(err, "CreateSessionEvents")
		}
		insertedEvents = append(insertedEvents, e)
	}

	return &insertedEvents, nil
}

const getSessionEventFromSessionID string = `
SELECT scheduled_event.* FROM scheduled_event
JOIN tutor_session ON tutor_session.event_id = scheduled_event.id
WHERE tutor_session.id = ?
	`

func (ss SessionStore) GetSessionEventFromSessionID(
	db rfrl.DB,
	ID int,
) (*rfrl.Event, error) {
	row := db.QueryRowx(getSessionEventFromSessionID, ID)

	var event rfrl.Event

	err := row.StructScan(&event)

	return &event, errors.Wrap(err, "GetSessionEventFromSessionID")
}

const deleteSessionEvents string = `
DELETE FROM scheduled_event
WHERE id in (?)
	`

func (ss SessionStore) DeleteSessionEvents(db rfrl.DB, eventIDs []int) error {
	query, args, err := sqlx.In(deleteSessionEvents, eventIDs)

	if err != nil {
		return errors.Wrap(err, "DeleteSessionEvents")
	}
	query = db.Rebind(query)

	_, err = db.Queryx(query, args...)

	return errors.Wrap(err, "DeleteSessionEvents")
}

const getSessionEventByID string = `
SELECT scheduled_event.* from scheduled_event
JOIN tutor_session ON tutor_session.event_id = scheduled_event.id
WHERE scheduled_event.id = $1 AND tutor_session.id = $2
	`

func (ss SessionStore) GetSessionEventByID(db rfrl.DB, sessionID int, ID int) (*rfrl.Event, error) {
	row := db.QueryRowx(getSessionEventByID, ID, sessionID)

	var event rfrl.Event

	err := row.StructScan(&event)

	return &event, errors.Wrap(err, "GetSessionEventByID")
}

func filterInclusiveDateRange(query sq.SelectBuilder, events *[]rfrl.Event) sq.SelectBuilder {
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

func (ss SessionStore) CheckAllClientSessionHasResponded(db rfrl.DB, id int) (bool, error) {
	var notAllClientsResponded null.Bool
	err := db.QueryRowx(checkAllClientSessionHasRespondedQuery, id).Scan(&notAllClientsResponded)

	if err != nil {
		return false, errors.Wrap(err, "CheckAllClientSessionHasResponded")
	}

	if !notAllClientsResponded.Valid {
		return false, errors.New("Unexpected invalid bool returned from database")
	}

	return !notAllClientsResponded.Bool, nil
}

const getSessionFromConferenceIDQuery string = `
SELECT * FROM tutor_session
WHERE conference_id = $1
`

func (ss SessionStore) GetSessionFromConferenceID(db rfrl.DB, conferenceID string) (*rfrl.Session, error) {
	var session rfrl.Session
	err := db.QueryRowx(getSessionFromConferenceIDQuery, conferenceID).StructScan(&session)

	return &session, errors.Wrap(err, "GetSessionFromConferenceID")
}

const checkClientsAttendedTutorSessionQuery string = `
SELECT count(client_id) FROM session_client
JOIN tutor_session ON tutor_session.id = session_client.session_id
WHERE client_id IN (?) 
AND can_attend = TRUE 
AND tutor_session.state = 'scheduled'
AND tutor_session.tutor_id = ?
GROUP BY session_id
	`

func (ss SessionStore) CheckClientsAttendedTutorSession(db rfrl.DB, tutorID string, clientIDs []string) (bool, error) {
	sql, args, err := sqlx.In(checkClientsAttendedTutorSessionQuery, clientIDs, tutorID)

	if err != nil {
		return false, errors.Wrap(err, "CheckClientsAttendedTutorSession")
	}

	sql = db.Rebind(sql)

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return false, errors.Wrap(err, "CheckClientsAttendedTutorSession")
	}

	for rows.Next() {
		var b int
		err := rows.Scan(&b)

		if err != nil {
			return false, errors.Wrap(err, "CheckClientsAttendedTutorSession")
		}

		if b == len(clientIDs) {
			return true, nil
		}
	}

	return false, nil
}
