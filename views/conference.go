package views

import (
	"net/http"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

type ConferenceView struct {
	SessionUseCase    tutorme.SessionUseCase
	ConferenceUseCase tutorme.ConferenceUseCase
}

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
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

func (cv *ConferenceView) ConnectToSessionClients(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	defer ws.Close()

	if err != nil {
		return nil
	}

	ws.SetReadLimit(tutorme.MaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(tutorme.PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(tutorme.PongWait)); return nil })

	conferenceID := c.Param("conferenceID")

	if _, err := uuid.Parse(conferenceID); err != nil {
		closeWebSocketWithError(ws, errors.Wrap(err, "Conference ID is not valid").Error())
		return nil
	}

	// Check if session exists
	_, err = cv.SessionUseCase.GetSessionFromConferenceID(conferenceID)

	if err != nil {
		closeWebSocketWithError(ws, err.Error())
		return nil
	}

	cv.ConferenceUseCase.Serve(ws, conferenceID)

	return nil
}
