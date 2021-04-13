package tutorme

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Session struct {
	ID              int       `db:"id" json:"id"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
	TutorID         string    `db:"tutor_id" json:"tutorId"`
	Tutor           Client    `json:"tutor"`
	UpdatedBy       string    `db:"updated_by" json:"updatedBy"`
	RoomID          string    `db:"room_id" json:"roomId"`
	Clients         []Client  `json:"clients"`
	State           string    `db:"state" json:"state"`
	TargetedEventID null.Int  `db:"event_id" json:"eventId"`
	CanAttend       null.Bool `db:"can_attend" json:"canAttend"`
	Event           Event     `json:"event"`
}

type Event struct {
	ID        int         `db:"id" json:"id"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time   `db:"updated_at" json:"updatedAt"`
	StartTime time.Time   `db:"start_time" json:"start"`
	EndTime   time.Time   `db:"end_time" json:"end"`
	Title     null.String `db:"title" json:"title"`
}

const (
	SCHEDULED string = "scheduled"
	PENDING   string = "pending"
)

// NewSession creates new Session
func NewSession(
	tutorID string,
	updatedBy string,
	roomID string,
	state string,
) *Session {
	return &Session{
		TutorID:   tutorID,
		UpdatedBy: updatedBy,
		RoomID:    roomID,
		State:     state,
	}
}

// NewEvent creates new Event
func NewEvent(
	start time.Time,
	end time.Time,
	title string,
) *Event {
	event := Event{
		StartTime: start,
		EndTime:   end,
		Title:     null.NewString(title, title != ""),
	}
	return &event
}

type SessionStore interface {
	GetSessionByClientID(db DB, clientID string, state string) (*[]Session, error)
	GetSessionByRoomID(db DB, clientID string, roomID string, state string) (*[]Session, error)
	GetSessionByID(db DB, clientID string, ID int) (*Session, error)
	GetSessionEventFromSessionID(db DB, ID int) (*Event, error)
	GetSessionByIDForUpdate(db DB, clientID string, ID int) (*Session, error)
	GetSessionEventByID(db DB, sessionID int, ID int) (*Event, error)
	CheckSessionsIsForClient(db DB, client string, sessionIDs []int) (bool, error)
	DeleteSession(db DB, ID int) error
	UpdateSession(db DB, ID int, by string, state string, EventID null.Int) (*Session, error)
	CreateSession(db DB, session *Session) (*Session, error)
	CreateSessionClients(db DB, sessionID int, clientIDs []string) (*[]Client, error)
	CreateSessionEvents(db DB, events []Event) (*[]Event, error)
	CreateClientSelectionOfEvent(db DB, sessionID int, clientID string, canAttend bool) error
	DeleteSessionEvents(db DB, eventIds []int) error
	CheckClientsAttendedTutorSession(db DB, tutorID string, clientIDs []string) (bool, error)
	CheckAllClientSessionHasResponded(db DB, ID int) (bool, error)
	GetSessionsEvent(db DB, sessionID []int) (map[int]*Event, error)
}

type SessionUseCase interface {
	CreateSession(tutorID string, updatedBy string, roomID string, clients []string, state string) (*Session, error)
	UpdateSession(ID int, updatedBy string, state string) (*Session, error)
	DeleteSession(clientID string, ID int) error
	GetSessionByID(clientID string, ID int) (*Session, error)
	GetSessionByRoomId(clientID string, roomID string, state string) (*[]Session, error)
	GetSessionByClientID(clientID string, state string) (*[]Session, error)
	GetSessionEventByID(sessionID int, ID int) (*Event, error)
	CreateSessionEvent(clientID string, ID int, event Event) (*Event, error)
	ClientActionOnSessionEvent(clientID string, sessionID int, canAttend bool) error
	GetSessionRelatedEvents(clientID string, sessionID int, start null.Time, end null.Time, state null.String) (*[]Event, error)
	CheckAllClientSessionHasResponded(ID int) (bool, error)
	CheckSessionsIsForClient(clientID string, sessionIDs []int) (bool, error)
	GetSessionsEvent(sessionIDs []int) (map[int]*Event, error)
}
