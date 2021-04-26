package publisher

import (
	"encoding/json"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/gommon/log"
)

type ConferencePublisher struct {
	Publisher tutorme.Publisher
}

func NewConferencePublisher(publisher tutorme.Publisher) *ConferencePublisher {
	return &ConferencePublisher{
		Publisher: publisher,
	}
}

type PublishingCode struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

func (cp *ConferencePublisher) PublishCode(codeID int, rawCode string, language string) error {
	code := tutorme.Code{
		ID:   codeID,
		Code: rawCode,
	}

	codeInJSON, err := json.Marshal(code)

	if err != nil {
		return err
	}

	topic := tutorme.CodeLanguageToTopic[language]

	log.Error(topic)

	err = cp.Publisher.Publish(topic, codeInJSON)

	return err
}
