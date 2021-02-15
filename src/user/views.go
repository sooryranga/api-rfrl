package user

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// UserPayload is the struct used to hold payload from /user
	UserPayload struct {
		ID        string `path:"id"`
		Email     string `json:"email" validate:"omitempty,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"profileImageURL"`
		About     string `json:"about"`
	}

	// EducationPayload is the struct used to create education
	EducationPaylod struct {
		Institution     string `json:"institution"`
		Degree          string `json:"degree"`
		FeildOfStudy    string `json:"field_of_study"`
		start           string `json:"start" validate:"datetime"`
		end             string `json:"end" validate:"omitempty, datetime"`
		InstitutionLogo string `json:"institution_logo"`
	}
)

// UserPayloadValidation validates user inputs
func UserPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(UserPayload)

	if payload.ID != "" {
		_, err := uuid.Parse(payload.ID)
		if err != nil {
			sl.ReportError(payload.Email, "id", "Id", "validUUID", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

func (h *Handler) CreateUser(c echo.Context) error {
	payload := UserPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
}

func (h *Handler) UpdateUser(c echo.Context) error {
	payload := UserPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
}
