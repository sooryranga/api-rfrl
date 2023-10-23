package tutorme

import (
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/guregu/null.v4"
)

// Client model
type Client struct {
	ID        string      `db:"id" json:"id" mapstructure:"id"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt" mapstructure:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updatedAt" mapstructure:"updated_at"`
	FirstName null.String `db:"first_name" json:"firstName" mapstructure:"first_name"`
	LastName  null.String `db:"last_name" json:"lastName" mapstructure:"last_name"`
	About     null.String `db:"about" json:"about" mapstructure:"about"`
	Email     null.String `db:"email" json:"email" mapstructure:"email"`
	Photo     null.String `db:"photo" json:"photo" mapstructure:"photo"`
}

// NewClient creates new client model struct
func NewClient(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) *Client {
	client := Client{
		FirstName: null.NewString(firstName, firstName != ""),
		LastName:  null.NewString(lastName, lastName != ""),
		About:     null.NewString(about, about != ""),
		Email:     null.NewString(email, email != ""),
		Photo:     null.NewString(photo, photo != ""),
	}
	log.Errorf("%v", client)
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

type ClientStore interface {
	GetClientFromID(db DB, id string) (*Client, error)
	CreateClient(db DB, client *Client) (*Client, error)
	UpdateClient(db DB, ID string, client *Client) (*Client, error)
	GetClientFromIDs(db DB, ID []string) (*[]Client, error)
}

type ClientUseCase interface {
	CreateClient(firstName string, lastName string, about string, email string, photo string) (*Client, error)
	UpdateClient(id string, firstName string, lastName string, about string, email string, photo string) (*Client, error)
	GetClient(id string) (*Client, error)
}
