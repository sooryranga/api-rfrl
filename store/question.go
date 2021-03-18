package store

import (
	"github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type QuestionStore struct{}

func NewQuestionStore() *QuestionStore {
	return &QuestionStore{}
}

const deleteTagsForQuestion string = `
DELETE FROM question_tags WHERE question_id = $1
`

func createTagsForQuestion(db tutorme.DB, questionID int, tagIDs []int) error {
	_, err := db.Queryx(deleteTagsForQuestion, questionID)

	if err != nil {
		return err
	}

	query := sq.Insert("question_tags").Columns("question_id", "tag_id")

	for i := 0; i < len(tagIDs); i++ {
		query = query.Values(questionID, tagIDs[i])
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	_, err = db.Queryx(sql, args...)

	return err
}

const getTagsSQL string = `
SELECT * FROM tags WHERE id IN (?)
`

func getTags(db tutorme.DB, tagIDs []int) (*[]tutorme.Tags, error) {
	query, args, err := sqlx.In(getTagsSQL, tagIDs)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)

	var tags []tutorme.Tags

	for rows.Next() {
		var tag tutorme.Tags
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, err
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

func getTagsForQuestion(db tutorme.DB, questionID int) (*[]tutorme.Tags, error) {
	rows, err := db.Queryx(getTagsForQuestionSQL, questionID)

	if err != nil {
		return nil, err
	}
	var tags []tutorme.Tags
	for rows.Next() {
		var tag tutorme.Tags
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, err
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

func getTagsForMultipleQuestions(db tutorme.DB, questionIDs []int) (*map[int][]tutorme.Tags, error) {
	questionIDToTags := make(map[int][]tutorme.Tags)

	if len(questionIDs) == 0 {
		return &questionIDToTags, nil
	}

	TagsIDtoQuestionID := make(map[int][]int)
	var tagIDs []int

	query, args, err := sqlx.In(getTagIDsForMultipleQuestionsSQL, questionIDs)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var result getTagsForMultipleQuestionsStruct
		err = rows.StructScan(&result)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	query = db.Rebind(query)

	rows, err = db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag tutorme.Tags
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, err
		}
		for _, questionID := range TagsIDtoQuestionID[tag.ID] {
			questionIDToTags[questionID] = append(questionIDToTags[questionID], tag)
		}
	}

	return &questionIDToTags, nil
}

func (qs *QuestionStore) CreateQuestion(db tutorme.DB, question tutorme.Question) (*tutorme.Question, error) {
	sql, args, err := sq.
		Insert("question").
		Columns("title", "body", "from_id", "images").
		Values(
			question.Title,
			question.Body,
			question.FromID,
			question.Images,
		).
		Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(sql, args...)
	var createdQuestion tutorme.Question
	err = row.StructScan(&createdQuestion)

	if len(question.TagIDs) > 0 {
		err = createTagsForQuestion(db, createdQuestion.ID, question.TagIDs)

		if err != nil {
			return nil, err
		}

		tags, err := getTags(db, question.TagIDs)

		if err != nil {
			return nil, err
		}

		createdQuestion.Tags = *tags
	}

	return &createdQuestion, err
}

func (qs *QuestionStore) UpdateQuestion(db tutorme.DB, clientID string, id int, title string, body string, images []string, tags []int) (*tutorme.Question, error) {
	query := sq.Update("question")
	if title != "" {
		query = query.Set("title", title)
	}

	if body != "" {
		query = query.Set("body", body)
	}

	if len(images) != 0 {
		query = query.Set("images", images)
	}

	sql, args, err := query.Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	row := db.QueryRowx(sql, args...)

	var updatedQuestion tutorme.Question
	err = row.StructScan(&updatedQuestion)

	if err != nil {
		return nil, err
	}

	if len(tags) > 0 {
		err = createTagsForQuestion(db, id, tags)

		if err != nil {
			return nil, err
		}
		tags, err := getTags(db, tags)
		if err != nil {
			return nil, err
		}
		updatedQuestion.Tags = *tags
	} else {
		tags, err := getTagsForQuestion(db, updatedQuestion.ID)
		if err != nil {
			return nil, err
		}
		updatedQuestion.Tags = *tags
	}

	return &updatedQuestion, err
}

const deleteQuestionSQL string = `
DELETE FROM question WHERE id = $1
`

func (qs *QuestionStore) DeleteQuestion(db tutorme.DB, id int) error {
	_, err := db.Queryx(deleteQuestionSQL, id)

	return err
}

const getQuestionFromIDSQL string = `
SELECT * FROM question WHERE id = $1
`

func (qs *QuestionStore) GetQuestion(db tutorme.DB, id int) (*tutorme.Question, error) {
	row := db.QueryRowx(getQuestionFromIDSQL, id)
	var question tutorme.Question

	err := row.StructScan(&question)

	if err != nil {
		return nil, err
	}

	tags, err := getTags(db, question.TagIDs)

	if err != nil {
		return nil, err
	}

	question.Tags = *tags

	return &question, nil
}

const getQuestionsSQL string = `
SELECT * FROM question
`

func (qs *QuestionStore) GetQuestions(db tutorme.DB) (*[]tutorme.Question, error) {
	rows, err := db.Queryx(getQuestionsSQL)

	if err != nil {
		return nil, err
	}

	idToQuestion := make(map[int]*tutorme.Question)
	var questionIds []int

	for rows.Next() {
		var question tutorme.Question
		err = rows.StructScan(&question)
		if err != nil {
			return nil, err
		}
		idToQuestion[question.ID] = &question
		questionIds = append(questionIds, question.ID)
	}

	questionTag, err := getTagsForMultipleQuestions(db, questionIds)

	if err != nil {
		return nil, err
	}

	questions := make([]tutorme.Question, 0, len(questionIds))

	for id, question := range idToQuestion {
		if tags, ok := (*questionTag)[id]; ok {
			question.Tags = tags
		}
		questions = append(questions, *question)
	}

	return &questions, err
}

const getQuestionsForClientSQL string = `
SELECT * FROM question WHERE from_id = $1
`

func (qs *QuestionStore) GetQuestionsForClient(db tutorme.DB, clientID string) (*[]tutorme.Question, error) {
	rows, err := db.Queryx(getQuestionsForClientSQL, clientID)

	if err != nil {
		return nil, err
	}

	idToQuestion := make(map[int]*tutorme.Question)
	var questionIds []int

	for rows.Next() {
		var question tutorme.Question
		err = rows.StructScan(&question)
		if err != nil {
			return nil, err
		}
		idToQuestion[question.ID] = &question
		questionIds = append(questionIds, question.ID)
	}

	questionTag, err := getTagsForMultipleQuestions(db, questionIds)

	if err != nil {
		return nil, err
	}

	questions := make([]tutorme.Question, 0, len(questionIds))

	for id, question := range idToQuestion {
		if tags, ok := (*questionTag)[id]; ok {
			question.Tags = tags
		}
		questions = append(questions, *question)
	}

	return &questions, err
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

func (qs *QuestionStore) ApplyToQuestion(db tutorme.DB, clientID string, id int) error {
	_, err := db.Queryx(insertQuestionApplicants, clientID, id)

	if err != nil {
		return err
	}

	_, err = db.Queryx(updateApplicantsOnQuestion, id)

	return err
}
