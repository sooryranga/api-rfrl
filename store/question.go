package store

import (
	"github.com/Arun4rangan/api-rfrl/rfrl"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type QuestionStore struct{}

func NewQuestionStore() *QuestionStore {
	return &QuestionStore{}
}

const deleteTagsForQuestion string = `
DELETE FROM question_tags WHERE question_id = $1
`

func createTagsForQuestion(db rfrl.DB, questionID int, tagIDs []int) error {
	_, err := db.Queryx(deleteTagsForQuestion, questionID)

	if err != nil {
		return errors.Wrap(err, "createTagsForQuestion")
	}

	query := sq.Insert("question_tags").Columns("question_id", "tag_id")

	for i := 0; i < len(tagIDs); i++ {
		query = query.Values(questionID, tagIDs[i])
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return errors.Wrap(err, "createTagsForQuestion")
	}

	_, err = db.Queryx(sql, args...)

	return errors.Wrap(err, "createTagsForQuestion")
}

const getTagsSQL string = `
SELECT * FROM tags WHERE id IN (?)
`

func getTags(db rfrl.DB, tagIDs []int) (*[]rfrl.Tags, error) {
	query, args, err := sqlx.In(getTagsSQL, tagIDs)
	if err != nil {
		return nil, errors.Wrap(err, "getTags")
	}
	query = db.Rebind(query)

	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "getTags")
	}

	tags := make([]rfrl.Tags, 0)

	for rows.Next() {
		var tag rfrl.Tags
		err = rows.StructScan(&tag)

		if err != nil {
			return nil, errors.Wrap(err, "getTags")
		}

		tags = append(tags, tag)
	}
	return &tags, nil
}

const getTagsForQuestionSQL string = `
SELECT tags.* FROM question_tags 
INNER JOIN tags ON question_tags.tag_id = tags.id
WHERE question_tags.question_id = $1
`

func getTagsForQuestion(db rfrl.DB, questionID int) (*[]rfrl.Tags, error) {
	rows, err := db.Queryx(getTagsForQuestionSQL, questionID)

	if err != nil {
		return nil, errors.Wrap(err, "getTagsForQuestion")
	}

	tags := make([]rfrl.Tags, 0)

	for rows.Next() {
		var tag rfrl.Tags
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, errors.Wrap(err, "getTagsForQuestion")
		}
		tags = append(tags, tag)
	}
	return &tags, nil
}

const getTagIDsForMultipleQuestionsSQL string = `
SELECT question_id, tag_id FROM question_tags WHERE question_id IN (?)
`

const getTagsFromIDsSQL string = `
SELECT * FROM tags WHERE id in (?)
`

type getTagsForMultipleQuestionsStruct struct {
	QuestionID int `db:"question_id"`
	TagID      int `db:"tag_id"`
}

func getTagsForMultipleQuestions(db rfrl.DB, questionIDs []int) (*map[int][]rfrl.Tags, error) {
	questionIDToTags := make(map[int][]rfrl.Tags)

	if len(questionIDs) == 0 {
		return &questionIDToTags, nil
	}

	TagsIDtoQuestionID := make(map[int][]int)
	var tagIDs []int

	query, args, err := sqlx.In(getTagIDsForMultipleQuestionsSQL, questionIDs)
	if err != nil {
		return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
	}
	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
	}

	for rows.Next() {
		var result getTagsForMultipleQuestionsStruct
		err = rows.StructScan(&result)
		if err != nil {
			return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
		}
		TagsIDtoQuestionID[result.TagID] = append(
			TagsIDtoQuestionID[result.TagID],
			result.QuestionID,
		)
		tagIDs = append(tagIDs, result.TagID)
	}

	if len(tagIDs) == 0 {
		return &questionIDToTags, nil
	}

	query, args, err = sqlx.In(getTagsFromIDsSQL, tagIDs)

	if err != nil {
		return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
	}
	query = db.Rebind(query)

	rows, err = db.Queryx(query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
	}

	for rows.Next() {
		var tag rfrl.Tags
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, errors.Wrap(err, "getTagsForMultipleQuestions")
		}
		for _, questionID := range TagsIDtoQuestionID[tag.ID] {
			questionIDToTags[questionID] = append(questionIDToTags[questionID], tag)
		}
	}

	return &questionIDToTags, nil
}

func (qs *QuestionStore) CreateQuestion(db rfrl.DB, question rfrl.Question) (*rfrl.Question, error) {
	sql, args, err := sq.
		Insert("question").
		Columns("title", "body", "from_id").
		Values(
			question.Title,
			question.Body,
			question.FromID,
		).
		Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "CreateQuestion")
	}

	row := db.QueryRowx(sql, args...)
	var createdQuestion rfrl.Question
	err = row.StructScan(&createdQuestion)

	if err != nil {
		return nil, errors.Wrap(err, "CreateQuestion")
	}

	if len(question.TagIDs) > 0 {
		err = createTagsForQuestion(db, createdQuestion.ID, question.TagIDs)

		if err != nil {
			return nil, errors.Wrap(err, "CreateQuestion")
		}

		tags, err := getTags(db, question.TagIDs)

		if err != nil {
			return nil, errors.Wrap(err, "CreateQuestion")
		}

		createdQuestion.Tags = *tags
	}

	return &createdQuestion, nil
}

func (qs *QuestionStore) UpdateQuestion(db rfrl.DB, clientID string, id int, title string, body string, tags []int, resolved null.Bool) (*rfrl.Question, error) {
	query := sq.Update("question")
	if title != "" {
		query = query.Set("title", title)
	}

	if body != "" {
		query = query.Set("body", body)
	}

	if resolved.Valid {
		query = query.Set("resolved", resolved)
	}

	sql, args, err := query.Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "UpdateQuestion")
	}

	row := db.QueryRowx(sql, args...)

	var updatedQuestion rfrl.Question
	err = row.StructScan(&updatedQuestion)

	if err != nil {
		return nil, errors.Wrap(err, "UpdateQuestion")
	}

	if len(tags) > 0 {
		err = createTagsForQuestion(db, id, tags)

		if err != nil {
			return nil, errors.Wrap(err, "UpdateQuestion")
		}
		tags, err := getTags(db, tags)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateQuestion")
		}
		updatedQuestion.Tags = *tags
	} else {
		tags, err := getTagsForQuestion(db, updatedQuestion.ID)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateQuestion")
		}
		updatedQuestion.Tags = *tags
	}

	return &updatedQuestion, nil
}

const deleteQuestionSQL string = `
DELETE FROM question WHERE id = $1
`

func (qs *QuestionStore) DeleteQuestion(db rfrl.DB, id int) error {
	_, err := db.Queryx(deleteQuestionSQL, id)

	return errors.Wrap(err, "DeleteQuestion")
}

const getQuestionFromIDSQL string = `
SELECT * FROM question WHERE id = $1
`

func (qs *QuestionStore) GetQuestion(db rfrl.DB, id int) (*rfrl.Question, error) {
	row := db.QueryRowx(getQuestionFromIDSQL, id)
	var question rfrl.Question

	err := row.StructScan(&question)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestion")
	}

	tags, err := getTags(db, question.TagIDs)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestion")
	}

	question.Tags = *tags

	return &question, nil
}

func (qs *QuestionStore) GetQuestions(db rfrl.DB, lastQuestion null.Int, resolved null.Bool) (*[]rfrl.Question, error) {
	query := sq.Select("*").From("question")

	if lastQuestion.Valid {
		query = query.Where(sq.Lt{"id": lastQuestion})
	}

	if resolved.Valid {
		query = query.Where(sq.Eq{"resolved": resolved})
	}

	sql, args, err := query.
		OrderBy("id DESC").
		Limit(50).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestions")
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestions")
	}

	idToQuestion := make(map[int]*rfrl.Question)
	var questionIds []int

	for rows.Next() {
		var question rfrl.Question
		err = rows.StructScan(&question)

		if err != nil {
			return nil, errors.Wrap(err, "GetQuestions")
		}

		idToQuestion[question.ID] = &question
		questionIds = append(questionIds, question.ID)
	}

	questionTag, err := getTagsForMultipleQuestions(db, questionIds)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestions")
	}

	questions := make([]rfrl.Question, 0, len(questionIds))

	for id, question := range idToQuestion {
		if tags, ok := (*questionTag)[id]; ok {
			question.Tags = tags
		}
		questions = append(questions, *question)
	}

	return &questions, nil
}

func (qs *QuestionStore) GetQuestionsForClient(db rfrl.DB, clientID string, resolved null.Bool) (*[]rfrl.Question, error) {
	query := sq.Select("*").From("question").Where(sq.Eq{"from_id": clientID})

	if resolved.Valid {
		query = query.Where(sq.Eq{"resolved": resolved})
	}

	sql, args, err := query.
		OrderBy("id DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestionsForClient")
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestionsForClient")
	}

	idToQuestion := make(map[int]*rfrl.Question)
	var questionIds []int

	for rows.Next() {
		var question rfrl.Question
		err = rows.StructScan(&question)
		if err != nil {
			return nil, errors.Wrap(err, "GetQuestionsForClient")
		}
		idToQuestion[question.ID] = &question
		questionIds = append(questionIds, question.ID)
	}

	questionTag, err := getTagsForMultipleQuestions(db, questionIds)

	if err != nil {
		return nil, errors.Wrap(err, "GetQuestionsForClient")
	}

	questions := make([]rfrl.Question, 0, len(questionIds))

	for id, question := range idToQuestion {
		if tags, ok := (*questionTag)[id]; ok {
			question.Tags = tags
		}
		questions = append(questions, *question)
	}

	return &questions, errors.Wrap(err, "GetQuestionsForClient")
}

const insertQuestionApplicants string = `
INSERT INTO question_applicants (applicant_id, question_id)
VALUES ($1, $2)
`

const updateApplicantsOnQuestion string = `
UPDATE question 
SET applicants = applicants + 1
WHERE id = $1
`

func (qs *QuestionStore) ApplyToQuestion(db rfrl.DB, clientID string, id int) error {
	_, err := db.Queryx(insertQuestionApplicants, clientID, id)

	if err != nil {
		return errors.Wrap(err, "ApplyToQuestion")
	}

	_, err = db.Queryx(updateApplicantsOnQuestion, id)

	return errors.Wrap(err, "ApplyToQuestion")
}
