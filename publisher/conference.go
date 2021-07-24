package publisher

import (
	"encoding/json"

	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/pkg/errors"
)

type ConferencePublisher struct {
	Publisher rfrl.Publisher
}

func NewConferencePublisher(publisher rfrl.Publisher) *ConferencePublisher {
	return &ConferencePublisher{
		Publisher: publisher,
	}
}

type PublishingCode struct {
	ID        int    `json:"id"`
	SessionID int    `json:"sessionId"`
	Code      string `json:"code"`
}

func (cp *ConferencePublisher) PublishCode(sessionID int, codeID int, rawCode string, language string) error {
	code := PublishingCode{
		ID:        codeID,
		SessionID: sessionID,
		Code:      rawCode,
	}

	codeInJSON, err := json.Marshal(code)

	if err != nil {
		return errors.Wrap(err, "PublishCode")
	}

	topic := rfrl.CodeLanguageToTopic[language]

	err = cp.Publisher.Publish(topic, codeInJSON)

	return errors.Wrap(err, "PublishCode")
}
