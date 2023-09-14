package user

import "github.com/jmoiron/sqlx"

// Store stores db instance
type Store struct {
	db *sqlx.DB
}

// NewStore creates auth store for querying
func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

const (
	getUserByID string = `
SELECT * FROM user
WHERE user.id = $1
	`
	insertUser string = `
INSERT INTO user (first_name, last_name, about, email, photo)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
	`
)

// GetUserFromID queries the database for user with id
func (s Store) GetUserFromID(id string) (*User, error) {
	var m User
	err := s.db.QueryRowx(getUserByID, id).StructScan(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateUser creates a new row for a user in the database
func (s Store) CreateUser(user *User) (int, error) {
	row := s.db.QueryRow(
		insertUser,
		user.FirstName,
		user.LastName,
		user.About,
		user.Email,
		user.Photo,
	)

	var id int

	err := row.Scan(&id)
	return id, err
}
