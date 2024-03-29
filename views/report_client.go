package views

import (
	"net/http"

	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	CreateReportPayload struct {
		Accused string `json:"accused" validate:"required,gte=0"`
		Cause   string `json:"cause"`
	}
	DeleteReportPayload struct {
		Reporter string `query:"reporter" validate:"required,gte=0"`
		Accused  string `query:"accused" validate:"required,gte=0"`
	}
)

type ReportClientView struct {
	ReportClientUseCase rfrl.ReportClientUseCase
}

func (r *ReportClientView) CreateReport(c echo.Context) error {
	payload := CreateReportPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateReport - Bind"))
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateReport - Validate"))
	}

	if _, err := uuid.Parse(payload.Accused); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid").SetInternal(errors.Wrap(err, "CreateReport - uuid.Parse"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	ReportClient := rfrl.NewReportClient(
		claims.ClientID,
		payload.Accused,
		payload.Cause,
	)

	err = r.ReportClientUseCase.CreateReport(ReportClient)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (r *ReportClientView) DeleteReport(c echo.Context) error {
	payload := DeleteReportPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "DeleteReport - Bind"))
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "DeleteReport - Validate"))
	}

	if _, err := uuid.Parse(payload.Accused); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid").SetInternal(errors.Wrap(err, "DeleteReport - uuid.Parse"))
	}

	if _, err := uuid.Parse(payload.Reporter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid").SetInternal(errors.Wrap(err, "DeleteReport - uuid.Parse"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if claims.Admin == false {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	report := rfrl.NewReportClient(
		payload.Reporter,
		payload.Accused,
		"",
	)

	err = r.ReportClientUseCase.DeleteReport(report)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (r *ReportClientView) GetReports(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	reports, err := r.ReportClientUseCase.GetReports()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, reports)
}
