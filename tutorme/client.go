package tutorme

import (
	"database/sql"
	"time"

	"github.com/labstack/gommon/log"
)

// Client model
type Client struct {
	ID        string         `db:"id" json:"id" mapstructure:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at" mapstructure:"updated_at"`
	FirstName sql.NullString `db:"first_name" json:"first_name" mapstructure:"first_name"`
	LastName  sql.NullString `db:"last_name" json:"last_name" mapstructure:"last_name"`
	About     sql.NullString `db:"about" json:"about" mapstructure:"about"`
	Email     sql.NullString `db:"email" json:"email" mapstructure:"email"`
	Photo     sql.NullString `db:"photo" json:"photo" mapstructure:"photo"`
}

// NewClient creates new client model struct
func NewClient(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) *Client {
	client := Client{}
	if firstName != "" {
		client.FirstName = sql.NullString{String: firstName, Valid: true}
	}
	if lastName != "" {
		client.LastName = sql.NullString{String: lastName, Valid: true}
	}
	if about != "" {
		client.About = sql.NullString{String: about, Valid: true}
	}
	if email != "" {
		client.Email = sql.NullString{String: email, Valid: true}
	}
	if photo != "" {
		client.Photo = sql.NullString{String: photo, Valid: true}
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
}

type ClientUseCase interface {
	CreateClient(firstName string, lastName string, about string, email string, photo string) (*Client, error)
	UpdateClient(id string, firstName string, lastName string, about string, email string, photo string) (*Client, error)
	GetClient(id string) (*Client, error)
}
