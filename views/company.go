package views

import (
	"net/http"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type (
	// SuggestCompanyPayload is the struct used to hold payload from /company/suggest
	SuggestCompanyPayload struct {
		Name        string `json:"name" validate:"gte=3"`
		EmailDomain string `json:"emailDomain" validate:"gte=3"`
	}

	UpdateCompanyViewPayload struct {
		Name     string      `json:"name" validate:"gte=3"`
		Photo    null.String `json:"photo"`
		Active   null.Bool   `json:"active"`
		Industry null.String `json:"industry"`
		About    null.String `json:"about"`
	}

	UpdateCompanyEmailViewPayload struct {
		Name        string `json:"name" validate:"gte=3"`
		EmailDomain string `json:"emailDomain" validate:"gte=2"`
		Active      bool   `json:"active"`
	}
)

type CompanyView struct {
	CompanyUseCase tutorme.CompanyUseCase
}

func (comv *CompanyView) SuggestCompanyView(c echo.Context) error {
	payload := SuggestCompanyPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := comv.CompanyUseCase.CreateSuggestion(payload.Name, payload.EmailDomain)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (comv *CompanyView) UpdateCompanyView(c echo.Context) error {
	payload := UpdateCompanyViewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if claims.Admin == false {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	company, err := comv.CompanyUseCase.UpdateCompany(
		payload.Name,
		payload.Photo,
		payload.Industry,
		payload.About,
		payload.Active,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, company)
}

func (comv *CompanyView) UpdateCompanyEmailView(c echo.Context) error {
	payload := UpdateCompanyEmailViewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if claims.Admin == false {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	err = comv.CompanyUseCase.UpdateCompanyEmail(
		payload.Name,
		payload.EmailDomain,
		payload.Active,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (comv *CompanyView) GetCompanies(c echo.Context) error {
	active := null.Bool{}
	err := active.UnmarshalText([]byte(c.QueryParam("active")))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if claims.Admin == false {
		active.Bool = false
	}

	companies, err := comv.CompanyUseCase.GetCompanies(active.Bool)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, companies)
}
