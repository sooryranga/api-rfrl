package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"gopkg.in/guregu/null.v4"
)

type FireStoreClient struct {
	UserRef *firestore.CollectionRef
	Client  *firestore.Client
}

func NewFireStore(client *firestore.Client) *FireStoreClient {
	return &FireStoreClient{
		UserRef: client.Collection("users"),
		Client:  client,
	}
}

type User struct {
	ID       string `firststore:"_id"`
	UserName string `firestore:"username"`
	Avatar   string `firestore:"avatar"` // in millions
}

func (fs *FireStoreClient) userName(firstName string, lastName string) string {
	return fmt.Sprintf("%s %s", firstName, lastName)
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

func (fs *FireStoreClient) UpdateClient(id string, photo null.String, firstName null.String, lastName null.String) error {
	ctx := context.Background()
	user := fs.UserRef.Doc(id)
	updates := make([]firestore.Update, 0)

	if photo.Valid {
		updates = append(updates, firestore.Update{Path: "avatar", Value: photo})
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
