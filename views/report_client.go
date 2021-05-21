package views

import (
	"net/http"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	ReportClientUseCase tutorme.ReportClientUseCase
}

func (r *ReportClientView) CreateReport(c echo.Context) error {
	payload := CreateReportPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := uuid.Parse(payload.Accused); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid")
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ReportClient := tutorme.NewReportClient(
		claims.ClientID,
		payload.Accused,
		payload.Cause,
	)

	err = r.ReportClientUseCase.CreateReport(ReportClient)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (r *ReportClientView) DeleteReport(c echo.Context) error {
	payload := DeleteReportPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := uuid.Parse(payload.Accused); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid")
	}

	if _, err := uuid.Parse(payload.Reporter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Accused is not a valid uuid")
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if claims.Admin == false {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	report := tutorme.NewReportClient(
		payload.Reporter,
		payload.Accused,
		"",
	)

	err = r.ReportClientUseCase.DeleteReport(report)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (r *ReportClientView) GetReports(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	reports, err := r.ReportClientUseCase.GetReports()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, reports)
}
