package user

import (
	"database/sql"
	"time"
)

// User model
type User struct {
	ID        int            `db:"id" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	FirstName sql.NullString `db:"first_name" json:"first_name"`
	LastName  sql.NullString `db:"last_name" json:"last_name"`
	About     sql.NullString `db:"about" json:"about"`
	Email     sql.NullString `db:"email" json:"email"`
	Photo     sql.NullString `db:"photo" json:"photo"`
}

func newUser(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) *User {
	user := User{}
	if firstName != "" {
		user.FirstName = sql.NullString{String: firstName, Valid: true}
	}
	if lastName != "" {
		user.LastName = sql.NullString{String: lastName, Valid: true}
	}
	if about != "" {
		user.About = sql.NullString{String: about, Valid: true}
	}
	if email != "" {
		user.Email = sql.NullString{String: email, Valid: true}
	}
	if photo != "" {
		user.Photo = sql.NullString{String: photo, Valid: true}
	}
	return &user
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
