package tutorme

import (
	"time"
)

type Tags struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"tag_name" json:"name"`
	About string `db:"about" json:"about"`
}

type Question struct {
	ID         int       `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt" mapstructure:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt" mapstructure:"updated_at"`
	Title      string    `db:"title" json:"title"`
	Body       string    `db:"body" json:"body"`
	TagIDs     []int     `db:"tag_ids" json:"-"`
	Tags       []Tags    `json:"tags"`
	FromID     string    `db:"from_id" json:"-"`
	From       Client    `json:"from"`
	Applicants int       `db:"applicants" json:"applicants"`
}

func NewQuestion(title string, body string, tags []int, fromClient string) Question {
	return Question{
		Title:  title,
		Body:   body,
		TagIDs: tags,
		FromID: fromClient,
	}
}

type QuestionUseCase interface {
	CreateQuestion(clientID string, title string, body string, tags []int) (*Question, error)
	UpdateQuestion(clientID string, id int, title string, body string, tags []int) (*Question, error)
	DeleteQuestion(clientID string, id int) error
	GetQuestion(id int) (*Question, error)
	GetQuestions() (*[]Question, error)
	GetQuestionsForClient(clientID string) (*[]Question, error)
	ApplyToQuestion(clientID string, id int) error
}

type QuestionStore interface {
	CreateQuestion(db DB, question Question) (*Question, error)
	UpdateQuestion(db DB, clientID string, id int, title string, body string, tags []int) (*Question, error)
	DeleteQuestion(db DB, id int) error
	GetQuestion(db DB, id int) (*Question, error)
	GetQuestions(db DB) (*[]Question, error)
	GetQuestionsForClient(db DB, clientID string) (*[]Question, error)
	ApplyToQuestion(db DB, clientID string, id int) error
}
