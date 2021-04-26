package tutorme

import (
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/guregu/null.v4"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 512
)

const (
	RUNNING     string = "running"
	NOT_RUNNING string = "not_running"
)

// Conference model
type Conference struct {
	SessionID    int       `db:"session_id" json:"session_id"`
	CodeState    string    `db:"code_state" json:"code_state"`
	LatestCodeID null.Int  `db:"latest_code" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

// Code model
type Code struct {
	ID     int         `db:"id" json:"-"`
	Code   string      `db:"code" json:"-"`
	Result null.String `db:"result" json:"result"`
}

var CodeLanguageToTopic = map[string]string{
	"javascript": JavascriptTopic,
	"python":     PythonTopic,
	"golang":     GoLangTopic,
}

type ConferenceUseCase interface {
	Serve(conn *websocket.Conn, conferenceID string)
	SubmitCode(sessionID int, code string, language string) error
}

type ConferenceStore interface {
	GetOrCreateConference(db DB, sessionID int) (*Conference, error)
	CreateNewCode(db DB, sessionID int, rawCode string) (*Code, error)
}

type ConferencePublisher interface {
	PublishCode(codeID int, code string, language string) error
}
