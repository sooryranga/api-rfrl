package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	SubmitCodePayload struct {
		SessionID int    `path:"sessionID"`
		Code      string `json:"code"`
		Language  string `json:"language"`
	}

	SetCodeResultPayload struct {
		SessionID int    `path:"sessionID"`
		Result    string `json:"result"`
		ID        int    `path:"ID"`
	}

	SubmitCodeResponse struct {
		CodeID int `json:"codeId"`
	}
)

type ConferenceView struct {
	SessionUseCase    rfrl.SessionUseCase
	ConferenceUseCase rfrl.ConferenceUseCase
}

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func closeWebSocketWithError(ws *websocket.Conn, errString string) {
	log.Errorf(errString)

	err := ws.WriteMessage(
		websocket.TextMessage,
		[]byte(errString),
	)

	if err != nil {
		log.Errorf(err.Error())
	}
}

func (cv *ConferenceView) ConnectToSessionSimplePeerClients(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	defer ws.Close()

	if err != nil {
		return nil
	}

	ws.SetReadLimit(rfrl.MaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(rfrl.PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(rfrl.PongWait)); return nil })

	conferenceID := c.Param("conferenceID")

	if _, err := uuid.Parse(conferenceID); err != nil {
		websocketError := rfrl.WebsocketError{Error: errors.Wrap(err, "Conference ID is not valid").Error()}
		rawError, _ := json.Marshal(websocketError)
		closeWebSocketWithError(ws, string(rawError))
		return nil
	}

	// Check if session exists
	_, err = cv.SessionUseCase.GetSessionFromConferenceID(conferenceID)

	if err != nil {
		closeWebSocketWithError(ws, err.Error())
		return nil
	}

	cv.ConferenceUseCase.Serve(ws, conferenceID, rfrl.SIMPLEPEER)

	return nil
}

func (cv *ConferenceView) ConnectToSessionYJSClients(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	defer ws.Close()

	if err != nil {
		return nil
	}

	ws.SetReadLimit(rfrl.MaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(rfrl.PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(rfrl.PongWait)); return nil })

	conferenceID := c.Param("conferenceID")

	if _, err := uuid.Parse(conferenceID); err != nil {
		websocketError := rfrl.WebsocketError{Error: errors.Wrap(err, "Conference ID is not valid").Error()}
		rawError, _ := json.Marshal(websocketError)
		closeWebSocketWithError(ws, string(rawError))
		return nil
	}

	// Check if session exists
	_, err = cv.SessionUseCase.GetSessionFromConferenceID(conferenceID)

	if err != nil {
		closeWebSocketWithError(ws, err.Error())
		return nil
	}

	cv.ConferenceUseCase.Serve(ws, conferenceID, rfrl.YJS)

	return nil
}

func (cv *ConferenceView) SubmitCode(c echo.Context) error {
	payload := SubmitCodePayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, ok := rfrl.CodeLanguageToTopic[payload.Language]

	if !ok {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Programming Language (%s) not supported", payload.Language),
		)
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	forClient, err := cv.SessionUseCase.CheckSessionsIsForClient(claims.ClientID, []int{payload.SessionID})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !forClient {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized for this session")
	}

	id, err := cv.ConferenceUseCase.SubmitCode(payload.SessionID, payload.Code, payload.Language)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, SubmitCodeResponse{
		CodeID: id,
	})
}

func (cv *ConferenceView) SetCodeResult(c echo.Context) error {
	payload := SetCodeResultPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := cv.ConferenceUseCase.SetCodeResult(
		payload.SessionID,
		payload.ID,
		payload.Result,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
