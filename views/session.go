package views

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	// DocumentPayload is the struct used to hold payload from /session
	SessionPayload struct {
		ID              int      `path:"id"`
		TutorID         string   `json:"tutorId" validate:"required,gte=0,lte=100"`
		RoomID          string   `json:"roomId" validate:"required, gte=0, lte=10"`
		ClientIDs       []string `json:"clientIds" validate:"required"`
		State           string   `json:"state" validate:"required"`
		TargetedEventID int      `json:"targetedEventId" validate:"omitempty"`
	}

	// SessionEventPayload is the struct used to hold payload from /session/:sessionId/event/:id
	SessionEventPayload struct {
		ID        int    `path:"id"`
		SessionID int    `path:"sessionID"`
		Start     string `json:"start" validate:"required, datetime, gte"`
		End       string `json:"end" validate:"required, datetime, gtfield=Start"`
		Title     string `json:"title" validate:"required, gte=0, lte=20"`
	}

	ClientSelectionOfSessionEventPayload struct {
		SessionID int   `path:'sessionID"`
		CanAttend *bool `json:"canAttend" validate:"required"`
	}

	GetSessionRelatedEventsEndpointPayload struct {
		SessionID int    `path:"sessionID"`
		StartTime string `query:"start" validate:"omitempty, datetime"`
		EndTime   string `query:"end" validate:"omitempty, datetime"`
		State     string `query:"state" validate:"omitempty,oneof= scheduled pending"`
	}

	GetSessionConferenceIDEndpointResponse struct {
		ConferenceID string `json:"conferenceID"`
	}
)

type SessionView struct {
	SessionUseCase tutorme.SessionUseCase
}

func (sv *SessionView) CreateSessionEndpoint(c echo.Context) error {
	payload := SessionPayload{
		State: tutorme.PENDING,
	}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Errorf(payload.State)
	session, err := sv.SessionUseCase.CreateSession(
		payload.TutorID,
		claims.ClientID,
		payload.RoomID,
		payload.ClientIDs,
		payload.State,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, session)
}

func (sv *SessionView) UpdateSessionEndpoint(c echo.Context) error {
	payload := SessionPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	forClient, err := sv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, []int{payload.ID})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !forClient {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to update this session")
	}

	session, err := sv.SessionUseCase.UpdateSession(payload.ID, claims.ClientID, payload.State)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if session.TargetedEventID.Valid {
		event, err := sv.SessionUseCase.GetSessionEventByID(session.ID, int(session.TargetedEventID.Int64))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		session.Event = event
	}

	return c.JSON(http.StatusOK, session)
}

func (sv *SessionView) DeleteSessionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not valid").Error())
	}
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = sv.SessionUseCase.DeleteSession(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (sv *SessionView) GetSessionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not a valid").Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	forClient, err := sv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, []int{ID})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !forClient {
		return echo.NewHTTPError(http.StatusBadRequest, "Session does not belong to client")
	}

	session, err := sv.SessionUseCase.GetSessionByID(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Errorj(log.JSON{"session": session})
	if session.TargetedEventID.Valid {
		event, err := sv.SessionUseCase.GetSessionEventByID(session.ID, int(session.TargetedEventID.Int64))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		session.Event = event
	}

	return c.JSON(http.StatusOK, session)
}

func (sv *SessionView) CreateSessionEventEndpoint(c echo.Context) error {
	payload := SessionEventPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	start, err := time.Parse(time.RFC3339, payload.Start)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	end, err := time.Parse(time.RFC3339, payload.End)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event := tutorme.NewEvent(start, end, payload.Title)

	event, err = sv.SessionUseCase.CreateSessionEvent(claims.ClientID, payload.SessionID, *event)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *event)
}

func (sv *SessionView) GetSessionConferenceIDEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("sessionID"))

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	forClient, err := sv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, []int{ID})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !forClient {
		return echo.NewHTTPError(http.StatusBadRequest, "Session does not belong to client")
	}

	session, err := sv.SessionUseCase.GetSessionByID(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !session.TargetedEventID.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "This session is not scheduled")
	}

	event, err := sv.SessionUseCase.GetSessionEventByID(session.ID, int(session.TargetedEventID.Int64))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	session.Event = event

	if session.Event.StartTime.Sub(time.Now()).Minutes() > 30 {
		return echo.NewHTTPError(
			http.StatusBadGateway, fmt.Sprintf(
				"Session event starts at : %s. Please try again closer to that time.",
				session.Event.StartTime.Format(time.RFC3339),
			),
		)
	}

	if session.Event.EndTime.Sub(time.Now()).Minutes() < -10 {
		return echo.NewHTTPError(
			http.StatusBadGateway,
			fmt.Sprintf(
				"Session event has ended at : %s.",
				session.Event.EndTime.Format(time.RFC3339),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		GetSessionConferenceIDEndpointResponse{ConferenceID: session.ConferenceID},
	)
}

func (sv *SessionView) GetSessionEventEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not a valid").Error())
	}

	sessionID, err := strconv.Atoi(c.Param("sessionID"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Session ID passed in is not a valid").Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	forClient, err := sv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, []int{sessionID})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !forClient {
		return echo.NewHTTPError(http.StatusBadRequest, "Session does not belong to client")
	}

	event, err := sv.SessionUseCase.GetSessionEventByID(sessionID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, event)
}

func (sv *SessionView) GetSessionsEndpoint(c echo.Context) error {
	roomID := c.QueryParam("roomId")
	state := c.QueryParam("state")

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessions := &([]tutorme.Session{})
	sessionIDs := make([]int, 0)

	if roomID != "" {
		sessions, err = sv.SessionUseCase.GetSessionByRoomId(claims.ClientID, roomID, state)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for i := 0; i < len(*sessions); i++ {
			sessionIDs = append(sessionIDs, (*sessions)[i].ID)
		}

		forClient, err := sv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, sessionIDs)

		if err != nil {
			return err
		}

		if !forClient {
			return errors.New("Room does not belong to client")
		}

	} else {
		sessions, err = sv.SessionUseCase.GetSessionByClientID(claims.ClientID, state)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for i := 0; i < len(*sessions); i++ {
			sessionIDs = append(sessionIDs, (*sessions)[i].ID)
		}
	}

	if len(sessionIDs) > 0 {
		sessionIDToEvent, err := sv.SessionUseCase.GetSessionsEvent(sessionIDs)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for i := 0; i < len(*sessions); i++ {
			session := (*sessions)[i]
			if event, ok := sessionIDToEvent[session.ID]; ok {
				(*sessions)[i].Event = event
			}
		}
	}

	return c.JSON(http.StatusOK, *sessions)
}

func (sv *SessionView) CreateClientActionOnSessionEvent(c echo.Context) error {
	payload := ClientSelectionOfSessionEventPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = sv.SessionUseCase.ClientActionOnSessionEvent(claims.ClientID, payload.SessionID, *payload.CanAttend)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	allClientsResponded, err := sv.SessionUseCase.CheckAllClientSessionHasResponded(payload.SessionID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if allClientsResponded {
		_, err = sv.SessionUseCase.UpdateSession(payload.SessionID, claims.ClientID, tutorme.SCHEDULED)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func (sv *SessionView) GetSessionRelatedEventsEndpoint(c echo.Context) error {
	payload := GetSessionRelatedEventsEndpointPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var start null.Time
	if payload.StartTime != "" {
		parsedStart, err := time.Parse(time.RFC3339, payload.StartTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		start = null.NewTime(parsedStart, true)
	}

	var end null.Time
	if payload.EndTime != "" {
		parsedEnd, err := time.Parse(time.RFC3339, payload.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		end = null.NewTime(parsedEnd, true)
	}

	var state null.String
	if payload.State == "" {
		state = null.NewString(tutorme.SCHEDULED, true)
	} else {
		state = null.NewString(payload.State, true)
	}

	events, err := sv.SessionUseCase.GetSessionRelatedEvents(claims.ClientID, payload.SessionID, start, end, state)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *events)
}
