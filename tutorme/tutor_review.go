package tutorme

import "time"

type TutorReview struct {
	ID         int       `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	TutorID    string    `db:"tutor_id" json:"tutorId"`
	FromID     string    `db:"from_id" json:"-"`
	FromClient Client    `json:"from"`
	Stars      int       `db:"stars" json:"stars"`
	Review     string    `db:"review" json:"review"`
}

type TutorReviewAggregate struct {
	TotalStars       int `db:"total_stars" json:"totalStars"`
	TotalReviewCount int `db:"total_review_count" json:"totalReviewCount"`
}

func NewTutorReview(TutorID string, Stars int, Review string) TutorReview {
	return TutorReview{
		TutorID: TutorID,
		Stars:   Stars,
		Review:  Review,
	}
}

type TutorReviewUseCase interface {
	CreateTutorReview(ClientID string, TutorID string, Stars int, Review string) (*TutorReview, error)
	UpdateTutorReview(ClientID string, ID int, Stars int, Reviews string) (*TutorReview, error)
	DeleteTutorReview(ClientID string, ID int) error
	GetTutorReview(ID int) (*TutorReview, error)
	GetTutorReviews(ClientID string) (*[]TutorReview, error)
	GetTutorReviewsAggregate(ClientID string) (*TutorReviewAggregate, error)
}

type TutorReviewStore interface {
	CreateTutorReview(db DB, clientID string, tutorReview *TutorReview) (*TutorReview, error)
	CheckTutorReviewForClient(db DB, clientID string, id int) (bool, error)
	UpdateTutorReview(db DB, id int, stars int, review string) (*TutorReview, error)
	DeleteTutorReview(db DB, id int) error
	GetTutorReview(db DB, id int) (*TutorReview, error)
	GetTutorReviews(db DB, tutorID string) (*[]TutorReview, error)
	GetTutorReviewsAggregate(db DB, clientID string) (*TutorReviewAggregate, error)
}
