package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"gopkg.in/guregu/null.v4"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
)

const (
	// SUBSCRIBE is a const
	SUBSCRIBE   string = "subscribe"
	UNSUBSCRIBE string = "unsubscribe"
	PUBLISH     string = "publish"
)

type WebSocketClient struct {
	hub *ConferenceHub

	ConferenceID string

	// The websocket connection.
	conn *websocket.Conn

	From string

	// Buffered channel of outbound messages.
	send chan []byte
}

type SignallingMessage struct {
	MessageType     string   `json:"type"`
	SubscribeTopics []string `json:"topics"`
	RawMessage      []byte
	ConferenceID    string
	From            string
	FromClient      *WebSocketClient
}

func (c *WebSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID, "from": c.From})
			return
		}

		var signalMessage SignallingMessage
		err = json.Unmarshal(message, &signalMessage)
		signalMessage.RawMessage = message
		signalMessage.ConferenceID = c.ConferenceID
		signalMessage.From = c.From
		signalMessage.FromClient = c

		if err != nil {
			log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID})
			return
		}

		switch signalMessage.MessageType {
		case SUBSCRIBE:
			if len(signalMessage.SubscribeTopics) != 1 && signalMessage.SubscribeTopics[0] != c.ConferenceID {
				log.Errorj(log.JSON{"signalMessage": signalMessage, "conferenceID": c.ConferenceID})
			}
		case UNSUBSCRIBE:
			return
		case PUBLISH:
			c.hub.broadcast <- signalMessage
		}
	}
}

// writePump pumps messages From the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes From this goroutine.
func (c *WebSocketClient) writePump() {

	ticker := time.NewTicker(tutorme.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(tutorme.WriteWait))
			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID})
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID})
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID})
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(tutorme.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorj(log.JSON{"error": err.Error(), "conferenceID": c.ConferenceID})
				return
			}
		}
	}
}

type ConferenceHub struct {
	// Registered clients.
	sessionConnectedClients map[string]map[*WebSocketClient]bool

	// Inbound messages From the clients.
	broadcast chan SignallingMessage

	// Register requests From the clients.
	register chan *WebSocketClient

	// Unregister requests From clients.
	unregister chan *WebSocketClient
}

func NewConferenceHub() *ConferenceHub {
	return &ConferenceHub{
		broadcast:               make(chan SignallingMessage),
		register:                make(chan *WebSocketClient),
		unregister:              make(chan *WebSocketClient),
		sessionConnectedClients: make(map[string]map[*WebSocketClient]bool),
	}
}

func (h ConferenceHub) Run() {
	for {
		select {

		case client := <-h.register:
			conferenceID := client.ConferenceID
			from := client.From
			id := fmt.Sprintf("%s-%s", from, conferenceID)
			conferenceClients, ok := h.sessionConnectedClients[id]

			if !ok {
				conferenceClients = make(map[*WebSocketClient]bool)
				h.sessionConnectedClients[id] = conferenceClients
			}
			conferenceClients[client] = true

		case client := <-h.unregister:
			conferenceID := client.ConferenceID
			from := client.From
			id := fmt.Sprintf("%s-%s", from, conferenceID)
			conferenceClients, ok := h.sessionConnectedClients[id]
			if !ok {
				continue
			}
			if _, ok := conferenceClients[client]; ok {
				delete(conferenceClients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			conferenceID := message.ConferenceID
			from := message.From
			id := fmt.Sprintf("%s-%s", from, conferenceID)
			conferenceClients, ok := h.sessionConnectedClients[id]

			if !ok {
				continue
			}

			for client := range conferenceClients {
				if client == message.FromClient {
					continue
				}
				select {
				case client.send <- message.RawMessage:
				default:
					log.Error("error")
					// close(client.send)
					// delete(conferenceClients, client)
				}
			}
		}
	}
}

type ConferenceUseCase struct {
	DB                  *sqlx.DB
	Hub                 *ConferenceHub
	ConferenceStore     tutorme.ConferenceStore
	ConferencePublisher tutorme.ConferencePublisher
	FireStore           tutorme.FireStoreClient
}

func NewConferenceUseCase(
	db *sqlx.DB,
	conferenceStore tutorme.ConferenceStore,
	hub *ConferenceHub,
	publisher tutorme.ConferencePublisher,
	fireStore tutorme.FireStoreClient,
) ConferenceUseCase {
	return ConferenceUseCase{
		DB:                  db,
		Hub:                 hub,
		ConferencePublisher: publisher,
		ConferenceStore:     conferenceStore,
		FireStore:           fireStore,
	}
}

func (cu ConferenceUseCase) Serve(conn *websocket.Conn, conferenceID string, from string) {
	client := &WebSocketClient{
		hub:          cu.Hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		ConferenceID: conferenceID,
		From:         from,
	}

	client.hub.register <- client

	go client.writePump()
	client.readPump()
}

func (cu ConferenceUseCase) SetCodeResult(sessionID int, ID int, result string) error {
	code := tutorme.Code{
		ID:     ID,
		Result: null.NewString(result, true),
	}

	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = cu.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	if *err != nil {
		return *err
	}

	*err = cu.FireStore.UpdateCode(sessionID, ID, result)

	if *err != nil {
		return *err
	}

	_, *err = cu.ConferenceStore.UpdateCode(tx, ID, code)

	if *err != nil {
		return *err
	}

	return nil
}

func (cu ConferenceUseCase) SubmitCode(sessionID int, rawCode string, language string) (int, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = cu.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	var conference *tutorme.Conference
	conference, *err = cu.ConferenceStore.GetOrCreateConference(cu.DB, sessionID)

	if *err != nil {
		return 0, *err
	}

	if conference.CodeState != tutorme.NOT_RUNNING && time.Now().Sub(conference.UpdatedAt).Minutes() < 1 {
		*err = errors.New("code is currently running right now")
		return 0, *err
	}

	var code *tutorme.Code
	code, *err = cu.ConferenceStore.CreateNewCode(cu.DB, sessionID, rawCode)

	if *err != nil {
		return 0, *err
	}

	*err = cu.ConferencePublisher.PublishCode(sessionID, code.ID, rawCode, language)

	return code.ID, *err
}
