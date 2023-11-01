package views

import (
	"net/http"
	"strconv"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
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
		SessionID int    `path:"sessionId"`
		Start     string `json:"start" validate:"required, datetime, gte"`
		End       string `json:"end" validate:"required, datetime, gtfield=Start"`
		Title     string `json:"title" validate:"required, gte=0, lte=20"`
	}

	ClientSelectionOfSessionEventPayload struct {
		SessionID int   `path:'sessionId"`
		CanAttend *bool `json:"canAttend" validate:"required"`
	}
)

type SessionView struct {
	SessionUseCase tutorme.SessionUseCase
}

func (sv *SessionView) CreateSessionEndpoint(c echo.Context) error {
	payload := SessionPayload{}

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
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
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
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
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
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
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
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, *event)
}

func (sv *SessionView) GetSessionEventEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not a valid").Error())
	}

	sessionID, err := strconv.Atoi(c.Param("sessionId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Session ID passed in is not a valid").Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := sv.SessionUseCase.GetSessionEventByID(claims.ClientID, sessionID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, event)
}

func (sv *SessionView) GetSessionsEndpoint(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	state := c.QueryParam("state")

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessions := &([]tutorme.Session{})

	log.Errorj(log.JSON{"sessions": sessions})
	if roomID != "" {
		sessions, err = sv.SessionUseCase.GetSessionByRoomId(claims.ClientID, roomID, state)
	} else {
		sessions, err = sv.SessionUseCase.GetSessionByClientID(claims.ClientID, state)
	}
	log.Errorj(log.JSON{"sessions": sessions})
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
