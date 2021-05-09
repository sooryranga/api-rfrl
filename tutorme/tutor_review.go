package tutorme

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type TutorReview struct {
	ID         int         `db:"id" json:"id"`
	CreatedAt  time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updatedAt"`
	TutorID    string      `db:"tutor_id" json:"tutorId"`
	FromID     string      `db:"from_id" json:"-"`
	FromClient Client      `json:"from"`
	Stars      null.Int    `db:"stars" json:"stars"`
	Review     null.String `db:"review" json:"review"`
	Headline   null.String `db:"headline" json:"headline"`
}

type TutorReviewAggregate struct {
	TotalStars       int `db:"total_stars" json:"totalStars"`
	TotalReviewCount int `db:"total_review_count" json:"totalReviewCount"`
}

type PendingTutorReview struct {
	TutorID        string `db:"tutor_id" json:"tutorId"`
	TutorFirstName string `db:"first_name" json:"firstName"`
	TutorLastName  string `db:"last_name" json:"lastName"`
}

func NewTutorReview(tutorID string, stars int, review string, headline string) TutorReview {
	return TutorReview{
		TutorID:  tutorID,
		Stars:    null.NewInt(int64(stars), stars != 0),
		Review:   null.NewString(review, review != ""),
		Headline: null.NewString(headline, headline != ""),
	}
}

type TutorReviewUseCase interface {
	CreateTutorReview(ClientID string, TutorID string, Stars int, Review string, Headline string) (*TutorReview, error)
	UpdateTutorReview(ClientID string, ID int, Stars int, Review string, Headline string) (*TutorReview, error)
	DeleteTutorReview(ClientID string, ID int) error
	GetTutorReview(ID int) (*TutorReview, error)
	GetTutorReviews(ClientID string) (*[]TutorReview, error)
	GetTutorReviewsAggregate(ClientID string) (*TutorReviewAggregate, error)
	GetPendingReviews(ClientID string) (*[]PendingTutorReview, error)
	CreatePendingReview(menteeID string, tutorID string) error
	DeletePendingReview(menteeID string, tutorID string) error
}

type TutorReviewStore interface {
	CreateTutorReview(db DB, clientID string, tutorReview *TutorReview) (*TutorReview, error)
	CheckTutorReviewForClient(db DB, clientID string, id int) (bool, error)
	UpdateTutorReview(db DB, tutorReview *TutorReview) (*TutorReview, error)
	DeleteTutorReview(db DB, id int) error
	GetTutorReview(db DB, id int) (*TutorReview, error)
	GetTutorReviews(db DB, tutorID string) (*[]TutorReview, error)
	GetTutorReviewsAggregate(db DB, clientID string) (*TutorReviewAggregate, error)
	GetPendingReviews(db DB, ClientID string) (*[]PendingTutorReview, error)
	CreatePendingReview(db DB, menteeID string, tutorID string) error
	DeletePendingReview(db DB, menteeID string, tutorID string) error
	CheckIfReviewAlreadyExists(db DB, menteeID string, tutorID string) (bool, error)
}
