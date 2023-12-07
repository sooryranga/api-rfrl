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
		Columns("tutor_id", "from_id", "stars", "review", "headline").
		Values(
			tutorReview.TutorID,
			ClientID,
			tutorReview.Stars.Int64,
			tutorReview.Review.String,
			tutorReview.Headline.String,
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

func (trs *TutorReviewStore) UpdateTutorReview(db tutorme.DB, review *tutorme.TutorReview) (*tutorme.TutorReview, error) {
	query := sq.Update("tutor_review")
	if review.Stars.Valid {
		query = query.Set("stars", review.Stars.Int64)
	}

	if review.Review.Valid {
		query = query.Set("review", review.Review.String)
	}

	if review.Headline.Valid {
		query = query.Set("headline", review.Headline.String)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": review.ID}).
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

const checkIfReviewAlreadyExistsQuery string = `
SELECT EXISTS (
	SELECT 1 FROM tutor_review 
	WHERE tutor_id = $1 AND from_id = $2
)
`

func (trs *TutorReviewStore) CheckIfReviewAlreadyExists(db tutorme.DB, menteeID string, tutorID string) (bool, error) {
	exists := false
	row := db.QueryRowx(checkIfReviewAlreadyExistsQuery, tutorID, menteeID)
	err := row.Scan(&exists)

	return exists, err
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

	tutorReviews := make([]tutorme.TutorReview, 0)

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
SELECT COALESCE(SUM(stars),0) as total_stars, COALESCE(COUNT(*),1) as total_review_count FROM tutor_review
WHERE tutor_review.tutor_id = $1
`

func (trs *TutorReviewStore) GetTutorReviewsAggregate(db tutorme.DB, tutorID string) (*tutorme.TutorReviewAggregate, error) {
	row := db.QueryRowx(getTutorReviewsAggregateByTutorID, tutorID)

	var aggregate tutorme.TutorReviewAggregate
	err := row.StructScan(&aggregate)

	return &aggregate, err
}

const getPendingReviewsQuery string = `
SELECT tutor_id, first_name, last_name 
FROM pending_tutor_review
JOIN client on pending_tutor_review.tutor_id = client.id
WHERE mentee_id = $1
`

func (trs *TutorReviewStore) GetPendingReviews(db tutorme.DB, menteeID string) (*[]tutorme.PendingTutorReview, error) {
	rows, err := db.Queryx(getPendingReviewsQuery, menteeID)

	pendingReviews := make([]tutorme.PendingTutorReview, 0)

	if err != nil {
		return &pendingReviews, err
	}

	for rows.Next() {
		var pendingReview tutorme.PendingTutorReview

		err = rows.StructScan(&pendingReview)

		if err != nil {
			return &pendingReviews, err
		}

		pendingReviews = append(pendingReviews, pendingReview)
	}

	return &pendingReviews, nil
}

const createPendingReviewQuery string = `
INSERT INTO pending_tutor_review (mentee_id, tutor_id)
VALUES ($1, $2)
`

func (trs *TutorReviewStore) CreatePendingReview(db tutorme.DB, menteeID string, tutorID string) error {
	_, err := db.Queryx(createPendingReviewQuery, menteeID, tutorID)

	return err
}

const deletePendingReviewQuery string = `
DELETE FROM pending_tutor_review
WHERE mentee_id = $1 AND tutor_id = $2
`

func (trs *TutorReviewStore) DeletePendingReview(db tutorme.DB, menteeID string, tutorID string) error {
	_, err := db.Queryx(deletePendingReviewQuery, menteeID, tutorID)

	return err
}
