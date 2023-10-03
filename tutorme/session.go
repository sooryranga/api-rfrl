package tutorme

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	TutorID   string    `db:"tutor_id" json:"tutor_id"`
	By        string    `db:"by" json:"by"`
	RoomID    string    `db:"room_id" json:"room_id"`
	Clients   []Client  `json:"clients"`
	State     string    `db:"state" json:"state"`
}

type Event struct {
	ID        int            `db:"int" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	Start     time.Time      `db:"start" json:"start"`
	End       time.Time      `db:"end" json:"json"`
	Title     sql.NullString `db:"title" json:"title"`
	SessionId int            `db:"session`
}

// NewSession creates new Session
func NewSession(
	tutorID string,
	by string,
	roomID string,
) *Session {
	return &Session{
		TutorID: tutorID,
		By:      by,
		RoomID:  roomID,
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
	UpdateSession(db DB, id int, by string, state string) (*Session, error)
	CreateSession(db DB, session *Session) (*Session, error)
	CreateSessionClients(db DB, sessionID int, clientIDs []string) (*[]Client, error)
	CreateSessionEvents(db DB, events []Event) (*[]Event, error)
	GetSessionEventFromClientIDs(db DB, clientIds []string, state bool) (*[]Event, error)
	DeleteSessionEvents(db DB, IDs []int) error
}
