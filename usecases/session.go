package usecases

import (
	"database/sql"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SessionUseCase holds all business related functions for session
type SessionUseCase struct {
	DB           *sqlx.DB
	SessionStore tutorme.SessionStore
}

func NewSessionUseCase(db sqlx.DB, sessionStore tutorme.SessionStore) *SessionUseCase {
	return &SessionUseCase{&db, sessionStore}
}

func (su *SessionUseCase) CreateSession(
	tutorID string,
	updatedBy string,
	roomID string,
	clients []string,
) (*tutorme.Session, error) {
	session := tutorme.NewSession(tutorID, updatedBy, roomID)
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = su.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	if *err != nil {
		return nil, *err
	}

	session, *err = su.SessionStore.CreateSession(tx, session)

	if *err != nil {
		return nil, *err
	}

	cl := &([]tutorme.Client{})

	cl, *err = su.SessionStore.CreateSessionClients(tx, session.ID, clients)

	if *err != nil {
		return nil, *err
	}

	session.Clients = *cl

	return session, *err
}

func (su SessionUseCase) UpdateSession(
	ID int,
	updatedBy string,
	state string,
) (*tutorme.Session, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = su.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	var session, updatedSession *tutorme.Session

	session, *err = su.SessionStore.GetSessionByIDForUpdate(tx, ID)

	if *err != nil {
		return nil, *err
	}

	//TODO: Add logic on what should be updated
	updatedSession, *err = su.SessionStore.UpdateSession(tx, ID, updatedBy, state, sql.NullInt64{Valid: false})

	if *err != nil {
		return session, *err
	}

	return updatedSession, nil
}

func (su SessionUseCase) GetSessionByID(clientID string, ID int) (*tutorme.Session, error) {
	session, err := su.SessionStore.GetSessionByID(su.DB, ID)

	if err != nil {
		return nil, err
	}

	forClient, err := su.SessionStore.CheckSessionsIsForClient(su.DB, clientID, []int{session.ID})

	if err != nil {
		return nil, err
	}

	if !forClient {
		return nil, errors.New("Session does not belong to client")
	}

	return session, nil
}

func (su SessionUseCase) GetSessionByRoomId(clientID string, roomID string, state string) (*[]tutorme.Session, error) {
	sessions, err := su.SessionStore.GetSessionByRoomID(su.DB, roomID, state)

	if err != nil {
		return nil, err
	}

	var sessionIDs []int
	for i := 0; i < len(*sessions); i++ {
		sessionIDs = append(sessionIDs, (*sessions)[i].ID)
	}

	forClient, err := su.SessionStore.CheckSessionsIsForClient(su.DB, clientID, sessionIDs)

	if err != nil {
		return nil, err
	}

	if !forClient {
		return nil, errors.New("Room does not belong to client")
	}

	return sessions, nil
}

func (su SessionUseCase) GetSessionByClientID(clientID string, state string) (*[]tutorme.Session, error) {
	return su.SessionStore.GetSessionByClientID(su.DB, clientID, state)
}

func (su SessionUseCase) DeleteSession(clientID string, ID int) error {
	session, err := su.SessionStore.GetSessionByID(su.DB, ID)

	if err != nil {
		return err
	}

	if session.TutorID != clientID {
		return errors.Errorf("Client %s cannot delete session", clientID)
	}

	if session.State == tutorme.SCHEDULED {
		return errors.Errorf("Cannot delete a scheduled session")
	}

	return su.SessionStore.DeleteSession(su.DB, ID)
}

func (su SessionUseCase) CreateSessionEvent(clientID string, ID int, event tutorme.Event) (*tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	var err = new(error)
	var tx *sqlx.Tx
	var session *tutorme.Session
	var isOverLapping bool
	insertedEvents := &([]tutorme.Event{})

	tx, *err = su.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	session, *err = su.SessionStore.GetSessionByID(tx, ID)

	if *err != nil {
		return nil, *err
	}

	if session.State == tutorme.SCHEDULED {
		*err = errors.New("Cannot change tutor event after scheduled")
		return nil, *err
	}

	forClient := false
	var clientIDs []string

	for i := 0; i < len(session.Clients); i++ {
		if session.Clients[i].ID == clientID {
			forClient = true
		}
		clientIDs = append(clientIDs, session.Clients[i].ID)
	}

	if !forClient {
		*err = errors.New("Session does not belong to client")
		return nil, *err
	}

	if session.TargetedEventID.Valid == true && session.TutorID != clientID {
		*err = errors.New("Only tutor can change scheduled event date")
		return nil, *err
	}

	isOverLapping, *err = su.SessionStore.CheckOverlapingEvents(tx, clientIDs, &[]tutorme.Event{event})

	if *err != nil {
		return nil, *err
	}

	if isOverLapping {
		*err = errors.New("Events overlap")
		return nil, *err
	}

	insertedEvents, *err = su.SessionStore.CreateSessionEvents(tx, []tutorme.Event{event})

	if *err != nil {
		return nil, *err
	}

	createdEvent := (*insertedEvents)[0]
	currentEvent := session.TargetedEventID

	_, *err = su.SessionStore.UpdateSession(tx, ID, clientID, "", sql.NullInt64{Valid: true, Int64: int64(createdEvent.ID)})

	*err = su.SessionStore.DeleteSessionEvents(tx, []int{int(currentEvent.Int64)})

	return &createdEvent, *err
}

func (su SessionUseCase) ClientActionOnSessionEvent(clientID string, sessionID int, canAttend bool) error {
	session, err := su.SessionStore.GetSessionByID(su.DB, sessionID)

	if err != nil {
		return err
	}

	forClient := false

	for i := 0; i < len(session.Clients); i++ {
		if session.Clients[i].ID == clientID {
			forClient = true
		}
	}

	if !forClient {
		return errors.New("Session does not belong to client")
	}

	return su.SessionStore.CreateClientSelectionOfEvent(su.DB, sessionID, clientID, canAttend)
}

func (su SessionUseCase) GetSessionEventByID(clientID string, sessionID int, ID int) (*tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	forClient, err := su.SessionStore.CheckSessionsIsForClient(su.DB, clientID, []int{sessionID})

	if err != nil {
		return nil, err
	}

	if !forClient {
		return nil, errors.New("Session does not belong to client")
	}

	return su.SessionStore.GetSessionEventByID(su.DB, sessionID, ID)
}
