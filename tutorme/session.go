package tutorme

import (
	sql "database/sql"
	"time"
)

type Session struct {
	ID              int           `db:"id" json:"id"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	TutorID         string        `db:"tutor_id" json:"tutor_id"`
	UpdatedBy       string        `db:"updated_by" json:"updated_by"`
	RoomID          string        `db:"room_id" json:"room_id"`
	Clients         []Client      `json:"clients"`
	State           string        `db:"state" json:"state"`
	TargetedEventID sql.NullInt64 `db:"event_id" json:"event_id"`
}

type Event struct {
	ID        int            `db:"int" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	Start     time.Time      `db:"start" json:"start"`
	End       time.Time      `db:"end" json:"end"`
	Title     sql.NullString `db:"title" json:"title"`
	SessionID int            `db:"session_id" json:"session_id"`
}

const (
	SCHEDULED string = "scheduled"
)

// NewSession creates new Session
func NewSession(
	tutorID string,
	updatedBy string,
	roomID string,
) *Session {
	return &Session{
		TutorID:   tutorID,
		UpdatedBy: updatedBy,
		RoomID:    roomID,
	}
}

// NewEvent creates new Event
func NewEvent(
	start time.Time,
	end time.Time,
	title string,
) *Event {
	event := Event{Start: start, End: end}
	if title != "" {
		event.Title = sql.NullString{String: title, Valid: false}
	}
	return &event
}

type SessionStore interface {
	GetSessionByClientID(db DB, clientID string, state string) (*[]Session, error)
	GetSessionByRoomID(db DB, roomID string, state string) (*[]Session, error)
	GetSessionByID(db DB, ID int) (*Session, error)
	GetSessionEventFromSessionID(db DB, clientIDs []string, state string) (*[]Event, error)
	GetSessionByIDForUpdate(db DB, ID int) (*Session, error)
	GetSessionEventByID(db DB, sessionID int, ID int) (*Event, error)
	CheckSessionsIsForClient(db DB, client string, sessionIDs []int) (bool, error)
	CheckOverlapingEvents(db DB, ID int, events *[]Event) (bool, error)
	DeleteSession(db DB, ID int) error
	UpdateSession(db DB, id int, by string, state string, EventID sql.NullInt64) (*Session, error)
	CreateSession(db DB, session *Session) (*Session, error)
	CreateSessionClients(db DB, sessionID int, clientIDs []string) (*[]Client, error)
	CreateSessionEvents(db DB, events []Event) (*[]Event, error)
	CreateClientSelectionOfEvent(db DB, sessionID int, clientID string, canAttend bool) error
	GetScheduledEventsFromClientIDs(db DB, clientIds []string, state bool) (*[]Event, error)
	DeleteSessionEvents(db DB, eventIds []int, sessionID int) error
	DeleteEventClient(db DB, sessionID int, clientID string, eventIDs []int) error
}

type SessionUseCase interface {
	CreateSession(tutorID string, by string, roomID string, clients []string) (Session, error)
	UpdateSession(id int, by string, state string) (Session, error)
	DeleteSession(clientID string, ID int) error
	GetSessionByID(clientID string, ID int) (*Session, error)
	GetSessionByRoomId(clientID string, roomID string, state string) (*[]Session, error)
	GetSessionByClientID(clientID string, state string) (*[]Session, error)
	GetSessionEventByID(clientID string, sessionID int, ID int) (*Event, error)
	CreateSessionEvents(clientID string, ID int, events *[]Event) (*[]Event, error)
	ClientActionOnSessionEvent(clientID string, sessionID int, eventIDs []int) error
}
