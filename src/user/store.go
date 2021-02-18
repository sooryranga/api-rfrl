package user

import (
	"github.com/Arun4rangan/api-tutorme/src/db"
	sq "github.com/Masterminds/squirrel"
)

const (
	getUserByID string = `
SELECT * FROM user
WHERE user.id = $1
	`
	insertUser string = `
INSERT INTO user (first_name, last_name, about, email, photo)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
	`
)

// GetUserFromID queries the database for user with id
func GetUserFromID(db db.DB, id string) (*User, error) {
	var m User
	err := db.QueryRowx(getUserByID, id).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateUser creates a new row for a user in the database
func CreateUser(db db.DB, user *User) (*User, error) {
	row := db.QueryRowx(
		insertUser,
		user.FirstName,
		user.LastName,
		user.About,
		user.Email,
		user.Photo,
	)

	var m User

	err := row.Scan(&m)
	return &m, err
}

// UpdateUser updates a user in the database
func UpdateUser(db db.DB, ID int, user *User) (*User, error) {
	updateQuery := sq.Update("user")
	if user.FirstName.Valid {
		updateQuery.Set("first_name", user.FirstName)
	}
	if user.LastName.Valid {
		updateQuery.Set("last_name", user.LastName)
	}
	if user.About.Valid {
		updateQuery.Set("about", user.About)
	}
	if user.Photo.Valid {
		updateQuery.Set("photo", user.Photo)
	}
	if user.Email.Valid {
		updateQuery.Set("email", user.Email)
	}

	sql, args, err := updateQuery.Where("id", ID).Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(
		sql,
		args...,
	)

	var m User

	err = row.Scan(&m)
	return &m, err
}
