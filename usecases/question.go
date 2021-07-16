package usecases

import (
	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type QuestionUseCase struct {
	DB            *sqlx.DB
	ClientStore   rfrl.ClientStore
	QuestionStore rfrl.QuestionStore
}

func NewQuestionUsesCase(
	db *sqlx.DB,
	clientStore rfrl.ClientStore,
	questionStore rfrl.QuestionStore,
) *QuestionUseCase {
	return &QuestionUseCase{db, clientStore, questionStore}
}

func (qu *QuestionUseCase) CreateQuestion(clientID string, title string, body string, tags []int) (*rfrl.Question, error) {
	question := rfrl.NewQuestion(title, body, tags, clientID)
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = qu.DB.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	var createdQuestion *rfrl.Question
	createdQuestion, *err = qu.QuestionStore.CreateQuestion(tx, question)

	if *err != nil {
		return nil, *err
	}

	var client *rfrl.Client
	client, *err = qu.ClientStore.GetClientFromID(qu.DB, createdQuestion.FromID)

	if *err != nil {
		return nil, *err
	}

	createdQuestion.From = *client
	return createdQuestion, nil
}

func (qu *QuestionUseCase) UpdateQuestion(clientID string, id int, title string, body string, tags []int, resolved null.Bool) (*rfrl.Question, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = qu.DB.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	var updatedQuestion *rfrl.Question
	updatedQuestion, *err = qu.QuestionStore.UpdateQuestion(tx, clientID, id, title, body, tags, resolved)

	if *err != nil {
		return nil, *err
	}

	var client *rfrl.Client
	client, *err = qu.ClientStore.GetClientFromID(qu.DB, updatedQuestion.FromID)

	if *err != nil {
		return nil, *err
	}

	updatedQuestion.From = *client
	return updatedQuestion, nil
}

func (qu *QuestionUseCase) DeleteQuestion(clientID string, id int) error {
	return qu.QuestionStore.DeleteQuestion(qu.DB, id)
}

func (qu *QuestionUseCase) GetQuestion(id int) (*rfrl.Question, error) {
	question, err := qu.QuestionStore.GetQuestion(qu.DB, id)

	if err != nil {
		return nil, err
	}

	client, err := qu.ClientStore.GetClientFromID(qu.DB, question.FromID)

	if err != nil {
		return nil, err
	}

	question.From = *client

	return question, err
}

func (qu *QuestionUseCase) GetQuestions(lastQuestion null.Int, resolved null.Bool) (*[]rfrl.Question, error) {
	questions, err := qu.QuestionStore.GetQuestions(qu.DB, lastQuestion, resolved)

	if len(*questions) == 0 {
		return questions, nil
	}

	if err != nil {
		return nil, err
	}

	var fromIDs []string
	for i := 0; i < len(*questions); i++ {
		question := (*questions)[i]
		fromIDs = append(fromIDs, question.FromID)
	}

	clients, err := qu.ClientStore.GetClientFromIDs(qu.DB, fromIDs)

	if err != nil {
		return nil, err
	}

	IDtoClient := make(map[string]*rfrl.Client)

	for i := 0; i < len(*clients); i++ {
		client := (*clients)[i]
		IDtoClient[client.ID] = &client
	}

	for i := 0; i < len(*questions); i++ {
		(*questions)[i].From = *IDtoClient[(*questions)[i].FromID]
	}

	return questions, nil
}

func (qu *QuestionUseCase) GetQuestionsForClient(clientID string, resolved null.Bool) (*[]rfrl.Question, error) {
	questions, err := qu.QuestionStore.GetQuestionsForClient(qu.DB, clientID, resolved)

	if err != nil {
		return nil, err
	}

	client, err := qu.ClientStore.GetClientFromID(qu.DB, clientID)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(*questions); i++ {
		question := (*questions)[i]
		question.From = *client
	}

	return questions, nil
}

func (qu *QuestionUseCase) ApplyToQuestion(clientID string, id int) error {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = qu.DB.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	*err = qu.QuestionStore.ApplyToQuestion(tx, clientID, id)

	return *err
}
