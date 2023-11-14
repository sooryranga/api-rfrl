package usecases

import (
	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

// SessionUseCase holds all business related functions for session
type SessionUseCase struct {
	DB           *sqlx.DB
	SessionStore tutorme.SessionStore
	ClientStore  tutorme.ClientStore
}

func NewSessionUseCase(db sqlx.DB, sessionStore tutorme.SessionStore, clientStore tutorme.ClientStore) *SessionUseCase {
	return &SessionUseCase{&db, sessionStore, clientStore}
}

func (su *SessionUseCase) CreateSession(
	tutorID string,
	updatedBy string,
	roomID string,
	clients []string,
	state string,
) (*tutorme.Session, error) {
	session := tutorme.NewSession(tutorID, updatedBy, roomID, state)
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

	log.Errorj(log.JSON{"session_clients": cl})

	session.Clients = *cl

	return session, *err
}

func (su SessionUseCase) CheckAllClientSessionHasResponded(
	ID int,
) (bool, error) {
	return su.SessionStore.CheckAllClientSessionHasResponded(su.DB, ID)
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

	session, *err = su.SessionStore.GetSessionByIDForUpdate(tx, updatedBy, ID)

	if *err != nil {
		return nil, *err
	}

	conferenceID := null.NewString("", false)
	if state == tutorme.SCHEDULED {
		conferenceID = null.NewString("gen_random_uuid()", true)
	}

	//TODO: Add logic on what should be updated
	updatedSession, *err = su.SessionStore.UpdateSession(tx, ID, updatedBy, state, null.NewInt(0, false), conferenceID)

	if *err != nil {
		return session, *err
	}

	updatedSession.CanAttend = session.CanAttend

	return updatedSession, nil
}

func (su SessionUseCase) GetSessionByID(clientID string, ID int) (*tutorme.Session, error) {
	session, err := su.SessionStore.GetSessionByID(su.DB, clientID, ID)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (su SessionUseCase) GetSessionByRoomId(clientID string, roomID string, state string) (*[]tutorme.Session, error) {
	sessions, err := su.SessionStore.GetSessionByRoomID(su.DB, clientID, roomID, state)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (su SessionUseCase) GetSessionByClientID(clientID string, state string) (*[]tutorme.Session, error) {
	return su.SessionStore.GetSessionByClientID(su.DB, clientID, state)
}

func canDeleteSession(clientID string, session tutorme.Session) error {
	log.Errorj(log.JSON{"clientid": clientID, "session": session})
	if session.State == tutorme.PENDING && session.UpdatedBy == clientID {
		return nil
	}

	if session.TutorID != clientID {
		return nil
	}

	return errors.Errorf("Client %s cannot delete session", clientID)
}

func (su SessionUseCase) DeleteSession(clientID string, ID int) error {
	session, err := su.SessionStore.GetSessionByID(su.DB, clientID, ID)

	if err != nil {
		return err
	}

	err = canDeleteSession(clientID, *session)

	if err != nil {
		return err
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

	session, *err = su.SessionStore.GetSessionByID(tx, clientID, ID)

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

	isOverLapping, *err = su.ClientStore.CheckOverlapingEventsByClientIDs(tx, clientIDs, &[]tutorme.Event{event})

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

	_, *err = su.SessionStore.UpdateSession(
		tx,
		ID,
		clientID,
		"",
		null.IntFrom(int64(createdEvent.ID)),
		null.NewString("", false),
	)

	if *err != nil {
		return nil, *err
	}

	if currentEvent.Valid {
		*err = su.SessionStore.DeleteSessionEvents(tx, []int{int(currentEvent.Int64)})

		if *err != nil {
			return nil, *err
		}
	}

	*err = su.SessionStore.CreateClientSelectionOfEvent(tx, ID, clientID, true)

	return &createdEvent, *err
}

func (su SessionUseCase) ClientActionOnSessionEvent(clientID string, sessionID int, canAttend bool) error {
	session, err := su.SessionStore.GetSessionByID(su.DB, clientID, sessionID)

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

	err = su.SessionStore.CreateClientSelectionOfEvent(su.DB, sessionID, clientID, canAttend)

	return err
}

func (su SessionUseCase) GetSessionRelatedEvents(
	clientID string,
	sessionID int,
	start null.Time,
	end null.Time,
	state null.String,
) (*[]tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	session, err := su.SessionStore.GetSessionByID(su.DB, clientID, sessionID)

	if err != nil {
		return nil, err
	}
	forClient := false
	clientIds := make([]string, len(session.Clients))

	for i := 0; i < len(session.Clients); i++ {
		if session.Clients[i].ID == clientID {
			forClient = true
		}
		clientIds[i] = session.Clients[i].ID
	}

	if !forClient {
		return nil, errors.New("Session does not belong to client")
	}

	return su.ClientStore.GetRelatedEventsByClientIDs(su.DB, clientIds, start, end, state)
}

func (su SessionUseCase) GetSessionEventByID(sessionID int, ID int) (*tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	return su.SessionStore.GetSessionEventByID(su.DB, sessionID, ID)
}

func (su SessionUseCase) GetSessionsEvent(sessionIDs []int) (map[int]*tutorme.Event, error) {
	// This will be a problem for the future because there is no guarantees that two parallel transaction will result in a unique event range
	return su.SessionStore.GetSessionsEvent(su.DB, sessionIDs)
}

func (su SessionUseCase) CheckSessionsIsForClient(clientID string, sessionIDs []int) (bool, error) {
	return su.SessionStore.CheckSessionsIsForClient(su.DB, clientID, sessionIDs)
}

func (su SessionUseCase) GetSessionFromConferenceID(conferenceID string) (*tutorme.Session, error) {
	return su.SessionStore.GetSessionFromConferenceID(su.DB, conferenceID)
}
