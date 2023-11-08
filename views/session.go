package views

import (
	"net/http"
	"strconv"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
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
		StartTime string `query:"startTime" validate:"omitempty, datetime"`
		EndTime   string `query:"endTime" validate:"omitempty, datetime"`
		State     string `query:"state" validate:"omitempty,oneof= scheduled pending"`
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

	session, err := sv.SessionUseCase.CreateSession(
		payload.TutorID,
		claims.ClientID,
		payload.RoomID,
		payload.ClientIDs,
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

	session, err := sv.SessionUseCase.UpdateSession(payload.ID, claims.ClientID, payload.State)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	session, err := sv.SessionUseCase.GetSessionByID(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	event, err := sv.SessionUseCase.GetSessionEventByID(claims.ClientID, sessionID, ID)

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

	if roomID != "" {
		sessions, err = sv.SessionUseCase.GetSessionByRoomId(claims.ClientID, roomID, state)
	} else {
		sessions, err = sv.SessionUseCase.GetSessionByClientID(claims.ClientID, state)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	var startTime null.Time
	if payload.StartTime != "" {
		start, err := time.Parse(time.RFC3339, payload.StartTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		startTime = null.NewTime(start, true)
	}

	var endTime null.Time
	if payload.EndTime != "" {
		end, err := time.Parse(time.RFC3339, payload.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		endTime = null.NewTime(end, true)
	}

	var state null.String
	if payload.State == "" {
		state = null.NewString(tutorme.SCHEDULED, true)
	} else {
		state = null.NewString(payload.State, true)
	}

	events, err := sv.SessionUseCase.GetSessionRelatedEvents(claims.ClientID, payload.SessionID, startTime, endTime, state)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *events)
}
