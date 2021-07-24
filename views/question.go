package views

import (
	"net/http"
	"strconv"

	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	QuestionPayload struct {
		ID       int       `path:"id"`
		Title    string    `json:"title" validate:"required,gte=0,lte=150"`
		Body     string    `json:body" validate:"required, gte=0"`
		Tags     []int     `json:"tags" validate:"numeric"`
		Resolved null.Bool `json:"resolved"`
	}
)

type QuestionView struct {
	QuestionUseCase rfrl.QuestionUseCase
}

func (qv *QuestionView) CreateQuestionEndpoint(c echo.Context) error {
	payload := QuestionPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateQuestionEndpoint - Bind"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	question, err := qv.QuestionUseCase.CreateQuestion(
		claims.ClientID,
		payload.Title,
		payload.Body,
		payload.Tags,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusCreated, question)
}

func (qv QuestionView) UpdateQuestionEndpoint(c echo.Context) error {
	payload := QuestionPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "UpdateQuestionEndpoint - Bind"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	question, err := qv.QuestionUseCase.UpdateQuestion(
		claims.ClientID,
		payload.ID,
		payload.Title,
		payload.Body,
		payload.Tags,
		payload.Resolved,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}
	return c.JSON(http.StatusOK, question)
}

func (qv QuestionView) DeleteQuestionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "DeleteQuestionEndpoint - Atoi"))
	}

	claims, err := rfrl.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	err = qv.QuestionUseCase.DeleteQuestion(claims.ClientID, ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}
	return c.NoContent(http.StatusOK)
}

func (qv QuestionView) GetQuestionEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetQuestionEndpoint - err"))
	}

	question, err := qv.QuestionUseCase.GetQuestion(ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, question)
}

func (qv QuestionView) GetQuestionsEndpoint(c echo.Context) error {
	lastQuestionParam := c.QueryParam("lastQuestion")

	var lastQuestion null.Int
	if lastQuestionParam != "" {
		i, err := strconv.Atoi(lastQuestionParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetQuestionsEndpoint - Atoi"))
		}
		lastQuestion = null.IntFrom(int64(i))
	}

	resolved := null.Bool{}
	err := resolved.UnmarshalText([]byte(c.QueryParam("withCompany")))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetQuestionsEndpoint - UnmarshalText"))
	}

	questions, err := qv.QuestionUseCase.GetQuestions(lastQuestion, resolved)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}
	return c.JSON(http.StatusOK, questions)
}

func (qv QuestionView) GetQuestionsFromClientEndpoint(c echo.Context) error {
	clientID := c.Param("id")

	if _, err := uuid.Parse(clientID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetQuestionsFromClientEndpoint - Parse"))
	}

	resolved := null.Bool{}
	err := resolved.UnmarshalText([]byte(c.QueryParam("withCompany")))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetQuestionsFromClientEndpoint - UnmarshalText"))
	}

	questions, err := qv.QuestionUseCase.GetQuestionsForClient(clientID, resolved)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, questions)
}

func (qv QuestionView) ApplyToQuestionEndpoint(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Question ID is not valid").SetInternal(errors.Wrap(err, "ApplyToQuestionEndpoint - Atoi"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	err = qv.QuestionUseCase.ApplyToQuestion(claims.ClientID, questionID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}
