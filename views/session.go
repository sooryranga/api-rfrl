package views

import (
	"net/http"
	"strconv"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// DocumentPayload is the struct used to hold payload from /session
	SessionPayload struct {
		ID              int      `path:"id"`
		TutorID         string   `json:"tutor_id" validate:"required,gte=0,lte=100"`
		RoomID          string   `json:"room_id" validate:"required, gte=0, lte=10"`
		ClientIDs       []string `json:"client_ids" validate:"required"`
		State           string   `json:"state" validate:"required"`
		TargetedEventID int      `json:"targeted_event_id" validate:"omitempty"`
	}

	// SessionEventPayload is the struct used to hold payload from /session/:session-id/event/:id
	SessionEventPayload struct {
		ID        int       `path:"id"`
		SessionID int       `path:"session-id"`
		Start     time.Time `json:"start" validate:"required,datetime"`
		End       time.Time `json:"end" validate:"required,datetime"`
		Title     string    `json:"title" validate:"required, gte=0, lte=20"`
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
	payload := &SessionPayload{}

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
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not valid"))
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
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not a valid"))
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
	payload := &SessionEventPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event := tutorme.NewEvent(payload.Start, payload.End, payload.Title)

	events, err := sv.SessionUseCase.CreateSessionEvents(claims.ClientID, payload.SessionID, &[]tutorme.Event{*event})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, (*events)[0])
}

func (sv *SessionView) GetSessionEventEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID passed in is not a valid"))
	}

	sessionID, err := strconv.Atoi(c.Param("session-id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Session ID passed in is not a valid"))
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

func (sv *SessionView) GetSessionsFromRoomEndpoint(c echo.Context) error {
	roomID := c.Param("id")
	state := c.QueryParam("state")

	if roomID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("ID is not passed in"))
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessions, err := sv.SessionUseCase.GetSessionByRoomId(claims.ClientID, roomID, state)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, sessions)
}
