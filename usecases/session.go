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

	tx, err := su.DB.Beginx()

	if err != nil {
		return nil, err
	}

	session, err = su.SessionStore.CreateSession(tx, session)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	cl, err := su.SessionStore.CreateSessionClients(tx, session.ID, clients)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	session.Clients = *cl

	return session, nil
}

func (su SessionUseCase) UpdateSession(
	id int,
	by string,
	state string,
) (*tutorme.Session, error) {
	tx, err := su.DB.Beginx()

	session, err := su.SessionStore.GetSessionByIDForUpdate(su.DB, id)

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return nil, errors.Wrap(rb, err.Error())
		}
		return nil, err
	}

	//TODO: Add logic on what should be updated

	updatedSession, err := su.SessionStore.UpdateSession(su.DB, id, by, state, sql.NullInt64{Valid: false})

	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			return session, errors.Wrap(rb, err.Error())
		}
		return session, err
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

	if session.TutorID == clientID {
		return errors.Errorf("Client %s cannot delete session", clientID)
	}

	return su.SessionStore.DeleteSession(su.DB, ID)
}

func (su SessionUseCase) CreateSessionEvents(clientID string, ID int, events *[]tutorme.Event) (*[]tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	forClient, err := su.SessionStore.CheckSessionsIsForClient(su.DB, clientID, []int{ID})

	if err != nil {
		return nil, err
	}

	if !forClient {
		return nil, errors.New("Session does not belong to client")
	}

	isOverLapping, err := su.SessionStore.CheckOverlapingEvents(su.DB, ID, events)

	if err != nil {
		return nil, err
	}

	if !isOverLapping {
		return nil, errors.New("Events overlap")
	}

	insertedEvents, err := su.SessionStore.CreateSessionEvents(su.DB, *events)

	return insertedEvents, err
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
