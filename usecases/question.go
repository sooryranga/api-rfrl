package usecases

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
)

type QuestionUseCase struct {
	DB            *sqlx.DB
	ClientStore   tutorme.ClientStore
	QuestionStore tutorme.QuestionStore
}

func NewQuestionUsesCase(
	db *sqlx.DB,
	clientStore tutorme.ClientStore,
	questionStore tutorme.QuestionStore,
) *QuestionUseCase {
	return &QuestionUseCase{db, clientStore, questionStore}
}

func (qu *QuestionUseCase) CreateQuestion(clientID string, title string, body string, images []string, tags []int) (*tutorme.Question, error) {
	question := tutorme.NewQuestion(title, body, images, tags, clientID)
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = qu.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	var createdQuestion *tutorme.Question
	createdQuestion, *err = qu.QuestionStore.CreateQuestion(tx, question)

	if *err != nil {
		return nil, *err
	}

	var client *tutorme.Client
	client, *err = qu.ClientStore.GetClientFromID(qu.DB, createdQuestion.FromID)

	if *err != nil {
		return nil, *err
	}

	createdQuestion.From = *client
	return createdQuestion, nil
}

func (qu *QuestionUseCase) UpdateQuestion(clientID string, id int, title string, body string, image []string, tags []int) (*tutorme.Question, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = qu.DB.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	var updatedQuestion *tutorme.Question
	updatedQuestion, *err = qu.QuestionStore.UpdateQuestion(tx, clientID, id, title, body, image, tags)

	if *err != nil {
		return nil, *err
	}

	var client *tutorme.Client
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

func (qu *QuestionUseCase) GetQuestion(id int) (*tutorme.Question, error) {
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

func (qu *QuestionUseCase) GetQuestions() (*[]tutorme.Question, error) {
	questions, err := qu.QuestionStore.GetQuestions(qu.DB)

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

	IDtoClient := make(map[string]*tutorme.Client)

	for i := 0; i < len(*clients); i++ {
		client := (*clients)[i]
		IDtoClient[client.ID] = &client
	}

	for i := 0; i < len(*questions); i++ {
		(*questions)[i].From = *IDtoClient[(*questions)[i].FromID]
	}

	return questions, nil
}

func (qu *QuestionUseCase) GetQuestionsForClient(clientID string) (*[]tutorme.Question, error) {
	questions, err := qu.QuestionStore.GetQuestionsForClient(qu.DB, clientID)

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

	defer tutorme.HandleTransactions(tx, err)

	*err = qu.QuestionStore.ApplyToQuestion(tx, clientID, id)

	return *err
}
