package user

import (
	"database/sql"
	"time"
)

// User model
type User struct {
	ID        int            `db:"id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	About     sql.NullString `db:"about"`
	Email     sql.NullString `db:"email"`
	Photo     sql.NullString `db:"photo"`
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
