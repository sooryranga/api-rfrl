package views

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	ClientPayload struct {
		ID                string    `path:"id"`
		Email             string    `json:"email" validate:"omitempty,email"`
		FirstName         string    `json:"firstName"`
		LastName          string    `json:"lastName"`
		Photo             string    `json:"photo"`
		About             string    `json:"about"`
		IsTutor           null.Bool `json:"isTutor"`
		LinkedInProfile   string    `json:"linkedInProfile"`
		GithubProfile     string    `json:"githubProfile"`
		YearsOfExperience null.Int  `json:"yearsOfExperience"`
		WorkTitle         string    `json:"workTitle"`
	}

	// EducationPayload is the struct used to create education
	EducationPayload struct {
		Institution  string `json:"institution"`
		Degree       string `json:"degree"`
		FieldOfStudy string `json:"fieldOfStudy"`
		StartYear    int    `json:"startYear"`
		EndYear      int    `json:"endYear"`
	}

	// VerifyEmailPayload is the struct used to verify email
	VerifyEmailPayload struct {
		Email    string `json:"email" validate:"email"`
		Type     string `json:"type" validate:"required,oneof= user work"`
		PassCode string `json:"passCode" validate:"omitempty,len=6,numeric"`
	}

	GetClientEventsEndpointEndpointPayload struct {
		ClientID  string `path:"clientID"`
		StartTime string `query:"start" validate:"omitempty, datetime"`
		EndTime   string `query:"end" validate:"omitempty, datetime"`
		State     string `query:"state" validate:"omitempty,oneof= scheduled pending"`
	}

	ClientCompanyReferralPayload struct {
		CompanyIds           []int `json:"companyIds"`
		IsLookingForReferral bool  `json:"isLookingForReferral"`
	}

	GetReferralCompanyResponse struct {
		CompanyIds []int `json:"companyIds"`
	}

	GetClientsEndpointPayload struct {
		FromCompanyIds           []int       `query:"fromCompanyIds"`
		IsTutor                  null.Bool   `query:"isTutor"`
		WantingReferralCompanyId null.Int    `query:"wantingReferralCompanyId"`
		LastTutor                null.String `query:"lastClient"`
		IncludeSelf              null.Bool   `query:"includeSelf"`
	}
)

// ClientPayloadValidation validates client inputs
func ClientPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(ClientPayload)

	if payload.ID != "" {
		_, err := uuid.Parse(payload.ID)
		if err != nil {
			sl.ReportError(payload.Email, "id", "Id", "validUUID", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

type ClientView struct {
	ClientUseCase rfrl.ClientUseCase
}

// CreateClientEndpoint view is an endpoint used to create client
func (cv *ClientView) CreateClientEndpoint(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You need to be an admin to create a user")
	}

	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateClientEndpoint - Bind"))
	}

	client, err := cv.ClientUseCase.CreateClient(
		payload.FirstName,
		payload.LastName,
		payload.About,
		payload.Email,
		payload.Photo,
		payload.IsTutor,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusCreated, client)
}

// UpdateClientEndpoint view is an endpoint uused to create client
func (cv *ClientView) UpdateClientEndpoint(c echo.Context) error {
	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(errors.Wrap(err, "UpdateClientEndpoint - Bind"))
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	if !claims.Admin && claims.ClientID != payload.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You cannot update this client")
	}

	params := rfrl.UpdateClientPayload{
		FirstName:         payload.FirstName,
		LastName:          payload.LastName,
		About:             payload.About,
		Email:             payload.Email,
		Photo:             payload.Photo,
		IsTutor:           payload.IsTutor,
		LinkedInProfile:   payload.LinkedInProfile,
		GithubProfile:     payload.GithubProfile,
		YearsOfExperience: payload.YearsOfExperience,
		WorkTitle:         payload.WorkTitle,
	}

	client, err := cv.ClientUseCase.UpdateClient(
		payload.ID,
		params,
	)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, "Client not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
		}
	}

	return c.JSON(http.StatusOK, client)
}

// GetClientEndpoint view is an endpoint uused to create client
func (cv *ClientView) GetClientEndpoint(c echo.Context) error {
	id := c.Param("id")

	client, err := cv.ClientUseCase.GetClient(id)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, "Client not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
		}
	}

	return c.JSON(http.StatusOK, client)
}

func (cv *ClientView) GetClientsEndpoint(c echo.Context) error {
	payload := GetClientsEndpointPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetClientsEndpoint - Bind"))
	}

	if payload.IsTutor.Valid && payload.IsTutor.Bool && payload.WantingReferralCompanyId.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot be looking for tutors wanting referrals")
	}

	if payload.IsTutor.Valid && !payload.IsTutor.Bool && len(payload.FromCompanyIds) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot be looking for clients from certain companies")
	}

	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	// Exclude Support
	excludeClients := []string{"af496484-7c7c-45a8-a409-96d61351f43a"}

	if !payload.IncludeSelf.Valid || !payload.IncludeSelf.Bool {
		excludeClients = append(excludeClients, claims.ClientID)
	}

	options := rfrl.GetClientsOptions{
		IsTutor:                  payload.IsTutor,
		CompanyIds:               payload.FromCompanyIds,
		WantingReferralCompanyId: payload.WantingReferralCompanyId,
		LastTutor:                payload.LastTutor,
		ExcludeClients:           excludeClients,
	}

	clients, err := cv.ClientUseCase.GetClients(options)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, clients)
}

func (cv *ClientView) VerifyEmail(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	payload := VerifyEmailPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "VerifyEmail - Bind"))
	}

	err = cv.ClientUseCase.CreateEmailVerification(
		claims.ClientID,
		payload.Email,
		payload.Type,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (cv *ClientView) VerifyEmailPassCode(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	payload := VerifyEmailPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "VerifyEmailPassCode - err"))
	}

	if len(payload.PassCode) != 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid passcode")
	}

	client, err := cv.ClientUseCase.VerifyEmail(
		claims.ClientID,
		payload.Email,
		payload.Type,
		payload.PassCode,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, client)
}

func (cv *ClientView) GetVerificationEmails(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	emailType := c.QueryParam("type")

	if emailType == "" {
		emailType = rfrl.WorkEmail
	} else if emailType != rfrl.WorkEmail && emailType != rfrl.UserEmail {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Only %s and %s are supported type for verification",
				rfrl.WorkEmail,
				rfrl.UserEmail,
			),
		)
	}

	email, err := cv.ClientUseCase.GetVerificationEmail(
		claims.ClientID,
		emailType,
	)

	if err != nil {
		if err.Error() == "Could not find verification email" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error()).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	type GetVerificationEmailsResponse struct {
		Email string `json:"email"`
	}

	return c.JSON(
		http.StatusOK,
		GetVerificationEmailsResponse{email},
	)

}

func (cv *ClientView) GetClientEventsEndpoint(c echo.Context) error {
	payload := GetClientEventsEndpointEndpointPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetClientEventsEndpoint - Bind"))
	}

	claims, err := rfrl.GetClaims(c)

	if !claims.Admin && claims.ClientID != payload.ClientID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to view this client")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	var start null.Time
	if payload.StartTime != "" {
		parsedStart, err := time.Parse(time.RFC3339, payload.StartTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetClientEventsEndpoint - time.Parse"))
		}

		start = null.NewTime(parsedStart, true)
	}

	var end null.Time
	if payload.EndTime != "" {
		parsedEnd, err := time.Parse(time.RFC3339, payload.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "GetClientEventsEndpoint - time.Parse"))
		}
		end = null.NewTime(parsedEnd, true)
	}

	var state null.String
	if payload.State == "" {
		state = null.NewString(rfrl.SCHEDULED, true)
	} else {
		state = null.NewString(payload.State, true)
	}

	events, err := cv.ClientUseCase.GetClientEvents(payload.ClientID, start, end, state)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, *events)
}

func (cv *ClientView) DeleteVerifyEmail(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	emailType := c.QueryParam("type")

	if emailType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Undefined type in query param")
	}

	if emailType != rfrl.WorkEmail && emailType != rfrl.UserEmail {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Only %s and %s are supported type for verification",
				rfrl.WorkEmail,
				rfrl.UserEmail,
			),
		)
	}

	err = cv.ClientUseCase.DeleteVerificationEmail(claims.ClientID, emailType)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (cv *ClientView) GetWantingReferralCompany(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	clientID := c.Param("clientID")

	if claims.ClientID != clientID && !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to edit this client")
	}

	companies, err := cv.ClientUseCase.GetClientWantingCompanyReferrals(
		clientID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error()).SetInternal(err)
	}

	response := GetReferralCompanyResponse{
		CompanyIds: companies,
	}
	return c.JSON(http.StatusOK, response)
}

func (cv *ClientView) CreateWantingReferralCompany(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	clientID := c.Param("clientID")

	if claims.ClientID != clientID && !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to edit this client")
	}

	payload := ClientCompanyReferralPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateWantingReferralCompany - Bind"))
	}

	err = cv.ClientUseCase.CreateClientWantingCompanyReferrals(
		clientID,
		payload.IsLookingForReferral,
		payload.CompanyIds,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (cv *ClientView) CreateEducation(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	clientID := c.Param("clientID")

	if claims.ClientID != clientID && !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to edit this client")
	}

	payload := EducationPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "CreateEducation - Bind"))
	}

	if payload.StartYear < 1950 || payload.StartYear > time.Now().Year()+5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong StartYear")
	}

	if payload.EndYear < 1950 || payload.StartYear > time.Now().Year()+10 {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong EndYear")
	}

	err = cv.ClientUseCase.CreateOrUpdateClientEducation(
		clientID,
		payload.Institution,
		payload.Degree,
		payload.FieldOfStudy,
		payload.StartYear,
		payload.EndYear,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	client, err := cv.ClientUseCase.GetClient(clientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}
	return c.JSON(http.StatusOK, client)
}
