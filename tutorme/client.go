package tutorme

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

const (
	WorkEmail = "work"
	UserEmail = "user"
)

// Client model
type Client struct {
	ID                string      `db:"id" json:"id"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updatedAt"`
	FirstName         null.String `db:"first_name" json:"firstName"`
	LastName          null.String `db:"last_name" json:"lastName"`
	About             null.String `db:"about" json:"about"`
	Email             null.String `db:"email" json:"email"`
	WorkEmail         null.String `db:"work_email" json:"work_email"`
	Photo             null.String `db:"photo" json:"photo"`
	IsTutor           null.Bool   `db:"is_tutor" json:"isTutor"`
	VerifiedWorkEmail null.Bool   `db:"verified_work_email" json:"verifiedWorkEmail"`
	VerifiedEmail     null.Bool   `db:"verified_email" json:"verifiedEmail"`
}

// NewClient creates new client model struct
func NewClient(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
	isTutor null.Bool,
) *Client {
	client := Client{
		FirstName: null.NewString(firstName, firstName != ""),
		LastName:  null.NewString(lastName, lastName != ""),
		About:     null.NewString(about, about != ""),
		Email:     null.NewString(email, email != ""),
		Photo:     null.NewString(photo, photo != ""),
		IsTutor:   isTutor,
	}

	return &client
}

// Education model
type Education struct {
	ID              int       `db:"id:`
	Institution     string    `db:"institution"`
	Degree          string    `db:"degree"`
	FieldOfStudy    string    `db:"field_of_study"`
	Start           time.Time `db:"start"`
	end             time.Time `db:"end"`
	InstitutionLogo string    `db:"institution_logo"`
}

type GetClientsOptions struct {
	IsTutor null.Bool
}

type ClientStore interface {
	GetClientFromID(db DB, id string) (*Client, error)
	CreateClient(db DB, client *Client) (*Client, error)
	UpdateClient(db DB, ID string, client *Client) (*Client, error)
	GetClientFromIDs(db DB, ID []string) (*[]Client, error)
	GetClients(db DB, options GetClientsOptions) (*[]Client, error)
	CreateEmailVerification(db DB, clientID string, email string, emailType string, passCode string) error
	VerifyEmail(db DB, clientID string, email string, emailType string, passCode string) error
	GetVerificationEmail(db DB, clientID string, emailType string) (string, error)
	DeleteVerificationEmail(db DB, clientID string, emailType string) error
}

type ClientUseCase interface {
	CreateClient(firstName string, lastName string, about string, email string, photo string, isTutor null.Bool) (*Client, error)
	UpdateClient(id string, firstName string, lastName string, about string, email string, photo string, isTutor null.Bool) (*Client, error)
	GetClient(id string) (*Client, error)
	GetClients(options GetClientsOptions) (*[]Client, error)
	CreateEmailVerification(clientID string, email string, emailType string) error
	VerifyEmail(clientID string, email string, emailType string, passCode string) (*Client, error)
	GetVerificationEmail(clientID string, emailType string) (string, error)
	DeleteVerificationEmail(clientID string, emailType string) error
}
