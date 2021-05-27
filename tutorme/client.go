package tutorme

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

const (
	WorkEmail = "work"
	UserEmail = "user"
)

// Education model
type Education struct {
	Institution     null.String `db:"institution" json:"institution"`
	Degree          null.String `db:"degree" json:"degree"`
	FieldOfStudy    null.String `db:"field_of_study" json:"fieldOfStudy"`
	StartYear       null.Int    `db:"start_year" json:"startYear"`
	EndYear         null.Int    `db:"end_year" json:"endYear"`
	InstitutionLogo null.String
}

func NewEducation(institution string, degree string, fieldOfStudy string, startYear int, endYear int) Education {
	return Education{
		Institution:  null.NewString(institution, institution != ""),
		FieldOfStudy: null.NewString(fieldOfStudy, fieldOfStudy != ""),
		Degree:       null.NewString(degree, degree != ""),
		StartYear:    null.NewInt(int64(startYear), startYear != 0),
		EndYear:      null.NewInt(int64(endYear), endYear != 0),
	}
}

// Client model
type Client struct {
	ID                   string      `db:"id" json:"id"`
	CreatedAt            time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time   `db:"updated_at" json:"updatedAt"`
	FirstName            null.String `db:"first_name" json:"firstName"`
	LastName             null.String `db:"last_name" json:"lastName"`
	About                null.String `db:"about" json:"about"`
	Email                null.String `db:"email" json:"email"`
	WorkEmail            null.String `db:"work_email" json:"workEmail"`
	CompanyID            null.Int    `db:"company_id" json:"companyID"`
	Photo                null.String `db:"photo" json:"photo"`
	IsTutor              null.Bool   `db:"is_tutor" json:"isTutor"`
	IsAdmin              null.Bool   `db:"is_admin" json:"-"`
	VerifiedWorkEmail    null.Bool   `db:"verified_work_email" json:"verifiedWorkEmail"`
	VerifiedEmail        null.Bool   `db:"verified_email" json:"verifiedEmail"`
	IsLookingForReferral null.Bool   `db:"is_looking_for_referral" json:"is_looking_for_referral"`
	Education
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
	GetRelatedEventsByClientIDs(db DB, clientIDs []string, start null.Time, end null.Time, state null.String) (*[]Event, error)
	CheckOverlapingEventsByClientIDs(db DB, clientIDs []string, events *[]Event) (bool, error)
	CreateOrUpdateClientEducation(db DB, clientID string, education Education) error
	CreateClientWantingCompanyReferrals(db DB, clientID string, companyIDs []int) error
	GetClientWantingCompanyReferrals(db DB, clientID string) ([]int, error)
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
	GetClientEvents(clientID string, start null.Time, end null.Time, state null.String) (*[]Event, error)
	CreateOrUpdateClientEducation(clientID string, institution string, degree string, fieldOfStudy string, startYear int, endYear int) error
	CreateClientWantingCompanyReferrals(clientID string, IsLookingForReferral bool, companyIds []int) error
	GetClientWantingCompanyReferrals(clientId string) ([]int, error)
}
