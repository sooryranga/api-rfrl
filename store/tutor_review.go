package store

import (
	"github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
)

type TutorReviewStore struct{}

func NewTutorReviewStore() *TutorReviewStore {
	return &TutorReviewStore{}
}

func (trs *TutorReviewStore) CreateTutorReview(db tutorme.DB, ClientID string, tutorReview *tutorme.TutorReview) (*tutorme.TutorReview, error) {
	sql, args, err := sq.Insert("tutor_review").
		Columns("tutor_id", "from_id", "stars", "review").
		Values(
			tutorReview.TutorID,
			ClientID,
			tutorReview.Stars,
			tutorReview.Review,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(sql, args...)

	var createdTutorReview tutorme.TutorReview
	err = row.StructScan(&createdTutorReview)

	return &createdTutorReview, err
}

const checkTutorReviewForClient string = `
SELECT count(*) FROM tutor_review 
WHERE client_id = $1 AND id = $2
	`

func (trs *TutorReviewStore) CheckTutorReviewForClient(db tutorme.DB, clientID string, id int) (bool, error) {
	row := db.QueryRowx(checkTutorReviewForClient, clientID, id)

	var reviewForClient int
	err := row.Scan(&reviewForClient)

	if err != nil {
		return false, err
	}
	return reviewForClient == 1, nil
}

func (trs *TutorReviewStore) UpdateTutorReview(db tutorme.DB, id int, stars int, review string) (*tutorme.TutorReview, error) {
	query := sq.Update("tutor_review")
	if stars != 0 {
		query = query.Set("stars", stars)
	}

	if review != "" {
		query = query.Set("review", review)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(sql, args...)

	var tutorReview tutorme.TutorReview
	err = row.StructScan(&tutorReview)

	return &tutorReview, err
}

const deleteTutorReview string = `
DELETE FROM tutor_review WHERE id = $1
	`

func (trs *TutorReviewStore) DeleteTutorReview(db tutorme.DB, ID int) error {
	_, err := db.Queryx(deleteTutorReview, ID)
	return err
}

const getTutorReview string = `
SELECT * FROM tutor_review 
WHERE tutor_review.id = $1
`

func (trs *TutorReviewStore) GetTutorReview(db tutorme.DB, id int) (*tutorme.TutorReview, error) {
	row := db.QueryRowx(getTutorReview, id)
	var tutorReview tutorme.TutorReview
	err := row.StructScan(&tutorReview)

	return &tutorReview, err
}

const getTutorReviewsByTutorID string = `
SELECT * FROM tutor_review
WHERE tutor_id = $1
`

func (trs *TutorReviewStore) GetTutorReviews(db tutorme.DB, tutorID string) (*[]tutorme.TutorReview, error) {
	rows, err := db.Queryx(getTutorReviewsByTutorID, tutorID)
	if err != nil {
		return nil, err
	}
	var tutorReviews []tutorme.TutorReview
	for rows.Next() {
		var tutorReview tutorme.TutorReview
		err = rows.StructScan(&tutorReview)
		if err != nil {
			return nil, err
		}
		tutorReviews = append(tutorReviews, tutorReview)
	}
	return &tutorReviews, nil
}

const getTutorReviewsAggregateByTutorID string = `
SELECT SUM(stars) as total_stars, COUNT(*) as total_review_count FROM tutor_review
WHERE tutor_review.tutor_id = $1
`

func (trs *TutorReviewStore) GetTutorReviewsAggregate(db tutorme.DB, tutorID string) (*tutorme.TutorReviewAggregate, error) {
	row := db.QueryRowx(getTutorReviewsAggregateByTutorID, tutorID)

	var aggregate tutorme.TutorReviewAggregate
	err := row.StructScan(&aggregate)

	return &aggregate, err
}