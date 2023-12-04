package views

import (
	"fmt"
	"net/http"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/guregu/null.v4"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	ClientPayload struct {
		ID        string    `path:"id"`
		Email     string    `json:"email" validate:"omitempty,email"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Photo     string    `json:"photo"`
		About     string    `json:"about"`
		IsTutor   null.Bool `json:"isTutor"`
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
	ClientUseCase tutorme.ClientUseCase
}

// CreateClientEndpoint view is an endpoint used to create client
func (cv *ClientView) CreateClientEndpoint(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot create a user")
	}

	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, client)
}

// UpdateClientEndpoint view is an endpoint uused to create client
func (cv *ClientView) UpdateClientEndpoint(c echo.Context) error {
	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if !claims.Admin && claims.ClientID != payload.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot update a user")
	}

	client, err := cv.ClientUseCase.UpdateClient(
		payload.ID,
		payload.FirstName,
		payload.LastName,
		payload.About,
		payload.Email,
		payload.Photo,
		payload.IsTutor,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

// GetClientEndpoint view is an endpoint uused to create client
func (cv *ClientView) GetClientEndpoint(c echo.Context) error {
	id := c.Param("id")

	client, err := cv.ClientUseCase.GetClient(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, client)

}

func (cv *ClientView) GetClientsEndpoint(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isTutor := null.Bool{}
	err = isTutor.UnmarshalText([]byte(c.QueryParam("is_tutor")))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	options := tutorme.GetClientsOptions{IsTutor: isTutor}

	if !claims.Admin {
		options.IsTutor = null.NewBool(true, true)
	}

	clients, err := cv.ClientUseCase.GetClients(options)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, clients)
}

func (cv *ClientView) VerifyEmail(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	payload := VerifyEmailPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = cv.ClientUseCase.CreateEmailVerification(
		claims.ClientID,
		payload.Email,
		payload.Type,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (cv *ClientView) VerifyEmailPassCode(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	payload := VerifyEmailPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(payload.PassCode) != 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid passcode")
	}

	log.Errorj(log.JSON{"VerifyEmailPayload": payload})
	client, err := cv.ClientUseCase.VerifyEmail(
		claims.ClientID,
		payload.Email,
		payload.Type,
		payload.PassCode,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

func (cv *ClientView) GetVerificationEmails(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	emailType := c.QueryParam("type")

	if emailType == "" {
		emailType = tutorme.WorkEmail
	} else if emailType != tutorme.WorkEmail && emailType != tutorme.UserEmail {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Only %s and %s are supported type for verification",
				tutorme.WorkEmail,
				tutorme.UserEmail,
			),
		)
	}

	email, err := cv.ClientUseCase.GetVerificationEmail(
		claims.ClientID,
		emailType,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := tutorme.GetClaims(c)

	if !claims.Admin && claims.ClientID != payload.ClientID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to view this client")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var start null.Time
	if payload.StartTime != "" {
		parsedStart, err := time.Parse(time.RFC3339, payload.StartTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		start = null.NewTime(parsedStart, true)
	}

	var end null.Time
	if payload.EndTime != "" {
		parsedEnd, err := time.Parse(time.RFC3339, payload.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		end = null.NewTime(parsedEnd, true)
	}

	var state null.String
	if payload.State == "" {
		state = null.NewString(tutorme.SCHEDULED, true)
	} else {
		state = null.NewString(payload.State, true)
	}

	events, err := cv.ClientUseCase.GetClientEvents(payload.ClientID, start, end, state)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, *events)
}

func (cv *ClientView) DeleteVerifyEmail(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	emailType := c.QueryParam("type")

	if emailType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Undefined type in query param")
	}

	if emailType != tutorme.WorkEmail && emailType != tutorme.UserEmail {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Only %s and %s are supported type for verification",
				tutorme.WorkEmail,
				tutorme.UserEmail,
			),
		)
	}

	err = cv.ClientUseCase.DeleteVerificationEmail(claims.ClientID, emailType)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (cv *ClientView) CreateEducation(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	clientID := c.Param("clientID")

	if claims.ClientID != clientID && !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to edit this client")
	}

	payload := EducationPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := cv.ClientUseCase.GetClient(clientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}
