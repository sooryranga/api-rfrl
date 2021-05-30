package views

import (
	"net/http"
	"strconv"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type (
	QuestionPayload struct {
		ID    int    `path:"id"`
		Title string `json:"title" validate:"required,gte=0,lte=150"`
		Body  string `json:body" validate:"required, gte=0"`
		Tags  []int  `json:"tags" validate:"numeric"`
	}
)

type QuestionView struct {
	QuestionUseCase tutorme.QuestionUseCase
}

func (qv *QuestionView) CreateQuestionEndpoint(c echo.Context) error {
	payload := QuestionPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	question, err := qv.QuestionUseCase.CreateQuestion(
		claims.ClientID,
		payload.Title,
		payload.Body,
		payload.Tags,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, question)
}

func (qv QuestionView) UpdateQuestionEndpoint(c echo.Context) error {
	payload := QuestionPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	question, err := qv.QuestionUseCase.UpdateQuestion(
		claims.ClientID,
		payload.ID,
		payload.Title,
		payload.Body,
		payload.Tags,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, question)
}

func (qv QuestionView) DeleteQuestionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = qv.QuestionUseCase.DeleteQuestion(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (qv QuestionView) GetQuestionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	question, err := qv.QuestionUseCase.GetQuestion(ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, question)
}

func (qv QuestionView) GetQuestionsEndpoint(c echo.Context) error {
	lastQuestionParam := c.QueryParam("lastQuestion")
	var lastQuestion null.Int
	if lastQuestionParam != "" {
		i, err := strconv.Atoi(lastQuestionParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}
		lastQuestion = null.IntFrom(int64(i))
	}

	questions, err := qv.QuestionUseCase.GetQuestions(lastQuestion)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, questions)
}

func (qv QuestionView) GetQuestionsFromClientEndpoint(c echo.Context) error {
	clientID := c.Param("id")

	if _, err := uuid.Parse(clientID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	questions, err := qv.QuestionUseCase.GetQuestionsForClient(clientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, questions)
}

func (qv QuestionView) ApplyToQuestionEndpoint(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Question ID is not valid")
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = qv.QuestionUseCase.ApplyToQuestion(claims.ClientID, questionID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
