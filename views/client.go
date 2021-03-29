package views

import (
	"net/http"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	ClientPayload struct {
		ID        string    `path:"id"`
		Email     string    `json:"email" validate:"omitempty,email"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Photo     string    `json:"profileImageURL"`
		About     string    `json:"about"`
		IsTutor   null.Bool `json:"is_tutor"`
	}

	// EducationPaylod is the struct used to create education
	EducationPaylod struct {
		Institution     string `json:"institution"`
		Degree          string `json:"degree"`
		FeildOfStudy    string `json:"fieldOfStudy"`
		start           string `json:"start" validate:"datetime"`
		end             string `json:"end" validate:"omitempty, datetime"`
		InstitutionLogo string `json:"institutionLogo"`
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
