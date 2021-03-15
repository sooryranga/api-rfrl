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
	ID        int            `db:"id" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	StartTime time.Time      `db:"start_time" json:"start_time"`
	EndTime   time.Time      `db:"end_time" json:"end_time"`
	Title     sql.NullString `db:"title" json:"title"`
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
	event := Event{StartTime: start, EndTime: end}
	if title != "" {
		event.Title = sql.NullString{String: title, Valid: false}
	}
	return &event
}

type SessionStore interface {
	GetSessionByClientID(db DB, clientID string, state string) (*[]Session, error)
	GetSessionByRoomID(db DB, roomID string, state string) (*[]Session, error)
	GetSessionByID(db DB, ID int) (*Session, error)
	GetSessionEventFromSessionID(db DB, ID int) (*Event, error)
	GetSessionByIDForUpdate(db DB, ID int) (*Session, error)
	GetSessionEventByID(db DB, sessionID int, ID int) (*Event, error)
	CheckSessionsIsForClient(db DB, client string, sessionIDs []int) (bool, error)
	CheckOverlapingEvents(db DB, clientIDs []string, events *[]Event) (bool, error)
	DeleteSession(db DB, ID int) error
	UpdateSession(db DB, ID int, by string, state string, EventID sql.NullInt64) (*Session, error)
	CreateSession(db DB, session *Session) (*Session, error)
	CreateSessionClients(db DB, sessionID int, clientIDs []string) (*[]Client, error)
	CreateSessionEvents(db DB, events []Event) (*[]Event, error)
	CreateClientSelectionOfEvent(db DB, sessionID int, clientID string, canAttend bool) error
	GetScheduledEventsFromClientIDs(db DB, clientIds []string, state string) (*[]Event, error)
	DeleteSessionEvents(db DB, eventIds []int) error
	CheckClientsAttendedTutorSession(db DB, tutorID string, clientIDs []string) (bool, error)
}

type SessionUseCase interface {
	CreateSession(tutorID string, updatedBy string, roomID string, clients []string) (*Session, error)
	UpdateSession(ID int, updatedBy string, state string) (*Session, error)
	DeleteSession(clientID string, ID int) error
	GetSessionByID(clientID string, ID int) (*Session, error)
	GetSessionByRoomId(clientID string, roomID string, state string) (*[]Session, error)
	GetSessionByClientID(clientID string, state string) (*[]Session, error)
	GetSessionEventByID(clientID string, sessionID int, ID int) (*Event, error)
	CreateSessionEvent(clientID string, ID int, event Event) (*Event, error)
	ClientActionOnSessionEvent(clientID string, sessionID int, canAttend bool) error
}
