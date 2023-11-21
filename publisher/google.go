package publisher

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/pubsub"
	pubsub "google.golang.org/api/pubsub/v1beta2"
)

type GooglePublisher struct {
	client   *pubsub.Client
	topicMap map[string]*pubsub.Topic
}

func NewGooglePublisher(projectID string) *GooglePublisher {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "project-id")
	if err != nil {
		log.Fatal(err)
	}

	return &GooglePublisher{
		client:   client,
		topicMap: make(map[string]*pubsub.Topic),
	}
}

func (p *GooglePublisher) CreateTopic(topicName string) error {
	ctx := context.Background()
	topic, err := p.client.CreateTopic(ctx, topicName)

	if err != nil {
		return err
	}

	p.topicMap[topicName] = topic

	return nil
}

func (p *GooglePublisher) Publish(topicName string, data []byte) error {
	ctx := context.Background()
	topic, ok := p.topicMap[topicName]

	if !ok {
		return errors.New("Topic not found")
	}

	publishResult := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
	})

	_, err := publishResult.Get(ctx)

	return err
}

func (p *GooglePublisher) Subscribe(topicName string, abort chan bool) (chan []byte, error) {
	messageChan := make(chan []byte, 1)

	ctx := context.Background()
	topic, ok := p.topicMap[topicName]

	if !ok {
		return nil, errors.New("Topic not found")
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
