package views

import (
	"net/http"
	"strconv"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	TutorReviewPayload struct {
		ID       int    `path:"id"`
		TutorID  string `json:"tutorId" validate:"required,gte=0,lte=100"`
		Stars    int    `json:"stars" validate:"required,numeric,gte=0,lte=10"`
		Review   string `json:"review"`
		Headline string `json:"headline"`
	}
)

type TutorReviewView struct {
	TutorReviewUseCase tutorme.TutorReviewUseCase
}

func (trv *TutorReviewView) CreateTutorReviewEndpoint(c echo.Context) error {
	payload := TutorReviewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if claims.ClientID == payload.TutorID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot create a review for yourself")
	}

	tutorReview, err := trv.TutorReviewUseCase.CreateTutorReview(
		claims.ClientID,
		payload.TutorID,
		payload.Stars,
		payload.Review,
		payload.Headline,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, tutorReview)
}

func (trv *TutorReviewView) UpdateTutorReviewEndpoint(c echo.Context) error {
	payload := TutorReviewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tutorReview, err := trv.TutorReviewUseCase.UpdateTutorReview(
		claims.ClientID,
		payload.ID,
		payload.Stars,
		payload.Review,
		payload.Headline,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tutorReview)
}

func (trv *TutorReviewView) DeleteTutorReviewEndpoint(c echo.Context) error {
	payload := TutorReviewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = trv.TutorReviewUseCase.DeleteTutorReview(claims.ClientID, payload.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (trv *TutorReviewView) GetTutorReviewEndpoint(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "ID pass in is not valid").Error())
	}

	tutorReview, err := trv.TutorReviewUseCase.GetTutorReview(ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tutorReview)
}

func (trv *TutorReviewView) GetTutorReviewsEndpoint(c echo.Context) error {
	clientID := c.Param("tutorID")
	// offest := c.QueryParam("offset")
	// pageLimit := c.QueryParam("page_limit")

	if clientID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tutor ID is not passed in")
	}

	tutorReviews, err := trv.TutorReviewUseCase.GetTutorReviews(clientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *tutorReviews)
}

func (trv *TutorReviewView) GetTutorReviewsAggregateEndpoint(c echo.Context) error {
	clientID := c.Param("tutorID")

	if clientID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Tutor ID is not passed in")
	}

	aggregateReview, err := trv.TutorReviewUseCase.GetTutorReviewsAggregate(clientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *aggregateReview)
}

func (trv *TutorReviewView) GetPendingReviewsEndpoint(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pendingReviewForClient, err := trv.TutorReviewUseCase.GetPendingReviews(
		claims.ClientID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, pendingReviewForClient)
}
