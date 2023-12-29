package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"gopkg.in/guregu/null.v4"
)

type FireStoreClient struct {
	ConferenceCodeRef *firestore.CollectionRef
	UserRef           *firestore.CollectionRef
	Client            *firestore.Client
	Auth              *auth.Client
}

func NewFireStore(client *firestore.Client, auth *auth.Client) *FireStoreClient {
	return &FireStoreClient{
		UserRef:           client.Collection("users"),
		ConferenceCodeRef: client.Collection("conferenceCode"),
		Client:            client,
		Auth:              auth,
	}
}

type User struct {
	ID       string `firststore:"_id"`
	UserName string `firestore:"username"`
	Avatar   string `firestore:"avatar"` // in millions
}

type Code struct {
	SessionID int       `firestore:"sessionId"`
	CodeID    int       `firestore:"codeId"`
	Timestamp time.Time `firestore:"timestamp"`
	Stdin     string    `firestore:"stdin"`
}

func (fs *FireStoreClient) userName(firstName string, lastName string) string {
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func (fs *FireStoreClient) CreateLoginToken(clientID string) (string, error) {
	ctx := context.Background()
	return fs.Auth.CustomToken(ctx, clientID)
}

func (fs *FireStoreClient) CreateClient(id string, photo string, firstName string, lastName string) error {
	ctx := context.Background()
	user := fs.UserRef.Doc(id)

	userName := fs.userName(firstName, lastName)
	_, err := user.Set(ctx, User{
		ID:       id,
		UserName: userName,
		Avatar:   photo,
	})
	return err
}

func (fs *FireStoreClient) CreateCode(sessionID int, codeID int) error {
	ctx := context.Background()
	code := fs.ConferenceCodeRef.Doc(fmt.Sprintf("%d-%d", sessionID, codeID))

	_, err := code.Set(ctx, Code{
		SessionID: sessionID,
		CodeID:    codeID,
		Timestamp: time.Now(),
		Stdin:     "run",
	})

	return err
}

func (fs *FireStoreClient) UpdateCode(sessionID int, codeID int, result string) error {
	ctx := context.Background()
	code := fs.ConferenceCodeRef.Doc(fmt.Sprintf("%d-%d", sessionID, codeID))

	_, err := code.Update(
		ctx,
		[]firestore.Update{{Path: "stdout", Value: result}},
	)

	return err
}

func (fs *FireStoreClient) UpdateClient(id string, photo null.String, firstName null.String, lastName null.String) error {
	ctx := context.Background()
	user := fs.UserRef.Doc(id)
	updates := make([]firestore.Update, 0)

	if photo.Valid {
		updates = append(updates, firestore.Update{Path: "avatar", Value: photo.String})
	}

	if firstName.Valid && lastName.Valid {
		userName := fs.userName(firstName.String, lastName.String)
		updates = append(updates, firestore.Update{Path: "username", Value: userName})
	}

	if len(updates) == 0 {
		return nil
	}
	_, err := user.Update(
		ctx,
		updates,
	)
	return err
}
