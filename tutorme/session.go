package tutorme

import (
	"database/sql"
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
	TargetedEventID sql.NullInt64 `db:"target_event_id" json:"target_event_id"`
}

type Event struct {
	ID        int            `db:"int" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	Start     time.Time      `db:"start" json:"start"`
	End       time.Time      `db:"end" json:"end"`
	Title     sql.NullString `db:"title" json:"title"`
	SessionId int            `db:"session_id" json:"session_id"`
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
		UpdatedBy: by,
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
	GetSessionByIDForUpdate(db DB, ID int) (*Session, error)
	CheckSessionsIsForClient(db DB, client string, sessionIDs []int) (bool, error)
	CheckOverlapingEvents(db DB, ID int, events *[]Event) (bool, error)
	DeleteSession(db DB, ID int) error
	UpdateSession(db DB, id int, by string, state string, targetEventID sql.NullInt64) (*Session, error)
	CreateSession(db DB, session *Session) (*Session, error)
	CreateSessionClients(db DB, sessionID int, clientIDs []string) (*[]Client, error)
	CreateSessionEvents(db DB, events []Event) (*[]Event, error)
	GetSessionEventFromClientIDs(db DB, clientIds []string, state bool) (*[]Event, error)
	DeleteSessionEvents(db DB, eventIds []int, sessionID int) error
	DeleteEventClient(db DB, sessionID int, clientID string, eventIDs []int) error
	CreateEventClient(db DB, sessionID int, clientID string, eventIDs []int) ([]CreateEventClientStoreResult, error)
}

type CreateEventClientStoreResult struct {
	Count   int `db:"count"`
	EventID int `db:"event_id"`
}

type SessionUseCase interface {
	CreateSession(tutorID string, by string, roomID string, clients []string) (Session, error)
	UpdateSession(id int, by string, state string) (Session, error)
	DeleteSession(clientID string, ID int) error
	GetSessionByID(clientID string, id int, roomID *string, state string) (*Session, error)
	GetSessionByRoomId(clientID string, roomID string, state string) (*[]Session, error)
	GetSessionByClientID(clientID string, state string) (*[]Session, error)
	CreateSessionEvents(clientID string, ID int, events *[]Event) (*[]Event, error)
	DeleteSessionEvents(clientID string, ID int, eventIDs []int) error
	SelectEvent(clientID string, sessionID int, eventIDs []int) error
	UnselectEvent(clientID string, sessionID int, eventIDs []int) error
}
