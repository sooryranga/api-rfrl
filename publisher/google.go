package publisher

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/gommon/log"
)

type GooglePublisher struct {
	client   *pubsub.Client
	topicMap map[string]*pubsub.Topic
}

func NewGooglePublisher(projectID string) (*GooglePublisher, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &GooglePublisher{
		client:   client,
		topicMap: make(map[string]*pubsub.Topic),
	}, nil
}

func (p *GooglePublisher) CreateTopic(topicName string) error {
	var topic *pubsub.Topic
	topic = p.client.Topic(topicName)

	p.topicMap[topicName] = topic

	return nil
}

func (p *GooglePublisher) Publish(topicName string, data []byte) error {
	ctx := context.Background()
	topic, ok := p.topicMap[topicName]

	if !ok {
		return errors.New("topic not found")
	}
	log.Error(data)
	publishResult := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	_, err := publishResult.Get(ctx)

	return err
}

func (p *GooglePublisher) Subscribe(topicName string, abort chan bool) (chan []byte, error) {
	messageChan := make(chan []byte, 1)

	ctx := context.Background()
	topic, ok := p.topicMap[topicName]

	if !ok {
		return nil, errors.New("topic not found")
	}

	sub, err := p.client.CreateSubscription(ctx, topicName, pubsub.SubscriptionConfig{
		Topic: topic,
	})

	if err != nil {
		return nil, err
	}

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		messageChan <- m.Data
		m.Ack() // Acknowledge that we've consumed the message.
	})

	return messageChan, err
}
