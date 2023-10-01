package tutorme

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	TutorID   string    `db:"tutor_id" json:"tutor_id"`
	By        string    `db:"by" json:"by"`
	RoomID    string    `db:"room_id" json:"room_id"`
}

type Event struct {
	ID        int            `db:"int" json:"id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
	Start     time.Time      `db:"start" json:"start"`
	End       time.Time      `db:"end" json:"json"`
	Title     sql.NullString `db:"title" json:"title"`
}
