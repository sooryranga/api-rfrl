package client

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// ClientPayload is the struct used to hold payload from /client
	ClientPayload struct {
		ID        string `path:"id"`
		Email     string `json:"email" validate:"omitempty,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"profileImageURL"`
		About     string `json:"about"`
	}

	// EducationPaylod is the struct used to create education
	EducationPaylod struct {
		Institution     string `json:"institution"`
		Degree          string `json:"degree"`
		FeildOfStudy    string `json:"field_of_study"`
		start           string `json:"start" validate:"datetime"`
		end             string `json:"end" validate:"omitempty, datetime"`
		InstitutionLogo string `json:"institution_logo"`
	}
)

// ClientPayloadValidation validates client inputs
func clientPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(ClientPayload)

	if payload.ID != "" {
		_, err := uuid.Parse(payload.ID)
		if err != nil {
			sl.ReportError(payload.Email, "id", "Id", "validUUID", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// CreateClientEndpoint view is an endpoint used to create client
func (h *Handler) CreateClientEndpoint(c echo.Context) error {
	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := h.createClient(
		payload.FirstName,
		payload.LastName,
		payload.About,
		payload.Email,
		payload.Photo,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, client)
}

// UpdateClientEndpoint view is an endpoint uused to create client
func (h *Handler) UpdateClientEndpoint(c echo.Context) error {
	payload := ClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := h.updateClient(
		payload.ID,
		payload.FirstName,
		payload.LastName,
		payload.About,
		payload.Email,
		payload.Photo,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

// GetClientEndpoint view is an endpoint uused to create client
func (h *Handler) GetClientEndpoint(c echo.Context) error {
	id := c.Param("id")

	client, err := h.getClient(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, client)

}