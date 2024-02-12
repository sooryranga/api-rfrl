package views

import (
	"net/http"
	"strconv"

	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type (
	CreateCompanyViewPayload struct {
		Name     string      `json:"name" validate:"gte=3"`
		Photo    null.String `json:"photo"`
		Active   null.Bool   `json:"active"`
		Industry null.String `json:"industry"`
		About    null.String `json:"about"`
	}

	UpdateCompanyViewPayload struct {
		ID       int         `path:"id" validate:"required"`
		Name     null.String `json:"name" validate:"gte=3"`
		Photo    null.String `json:"photo"`
		Active   null.Bool   `json:"active"`
		Industry null.String `json:"industry"`
		About    null.String `json:"about"`
	}

	UpdateCompanyEmailViewPayload struct {
		CompanyName string `json:"companyName"`
		EmailDomain string `json:"emailDomain" validate:"gte=2"`
		Active      *bool  `json:"active" validate:"required"`
	}
)

type CompanyView struct {
	CompanyUseCase rfrl.CompanyUseCase
}

func (comv *CompanyView) GetCompany(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "name cannot be empty")
	}

	company, err := comv.CompanyUseCase.GetCompany(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, company)
}

func (comv *CompanyView) CreateCompanyView(c echo.Context) error {
	payload := CreateCompanyViewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	company, err := comv.CompanyUseCase.CreateCompany(
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

func (comv *CompanyView) UpdateCompanyView(c echo.Context) error {
	payload := UpdateCompanyViewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	company, err := comv.CompanyUseCase.UpdateCompany(
		payload.ID,
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

func (comv *CompanyView) GetCompanyEmailsView(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	withCompany := null.Bool{}
	err = withCompany.UnmarshalText([]byte(c.QueryParam("withCompany")))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	companyEmails, err := comv.CompanyUseCase.GetCompanyEmails(withCompany)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, companyEmails)
}

func (comv *CompanyView) UpdateCompanyEmailView(c echo.Context) error {
	payload := UpdateCompanyEmailViewPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	err = comv.CompanyUseCase.UpdateCompanyEmail(
		null.NewString(payload.CompanyName, payload.CompanyName != ""),
		payload.EmailDomain,
		*payload.Active,
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

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin && active.Valid && !active.Bool {
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot filter for non active companies")
	}

	companies, err := comv.CompanyUseCase.GetCompanies(active.Bool)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, companies)
}

func (comv *CompanyView) GetCompanyEmailView(c echo.Context) error {
	companyEmail := c.Param("companyEmail")

	if companyEmail == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "you need to pass in company email")
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized for this view")
	}

	company, err := comv.CompanyUseCase.GetCompanyEmail(companyEmail)

	return c.JSON(http.StatusOK, company)
}
