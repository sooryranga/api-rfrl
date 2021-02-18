package auth

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// SignUpPayload is the struct used to hold payload from /signup
	SignUpPayload struct {
		Email     string `json:"email" validate:"omitempty,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Token     string `json:"token"`
		Photo     string `json:"profileImageURL"`
		About     string `json:"about"`
		Password  string `json:"password" validate:"omitempty,gte=10"`
		Type      string `json:"type" validate:"required,oneof= google linkedin email"`
	}
)

// SignUpPayloadValidation validates user inputs
func SignUpPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(SignUpPayload)

	switch payload.Type {
	case GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// Signup endpoint
func (h *Handler) Signup(c echo.Context) error {
	payload := SignUpPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var token string
	var err error

	switch payload.Type {
	case GOOGLE:
		token, err = h.signupGoogle(payload.Token)
	case LINKEDIN:
		token, err = h.signupLinkedIn(payload.Token)
	case EMAIL:
		token, err = h.signupEmail(payload.Email, payload.Password)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		c.Logger().Error(err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("User is already signed up with %s", payload.Type),
				)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

type (
	// LoginPayload is the struct used to hold payload from /login
	LoginPayload struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Token    string `json:"token"`
		Password string `json:"password" validate:"omitempty,gte=10"`
		Type     string `json:"type" validate:"required,oneof= google linkedin email"`
	}
)

// LoginPayloadValidation validates user inputs
func LoginPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(LoginPayload)

	switch payload.Type {
	case GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// Login is used to login an user
func (h *Handler) Login(c echo.Context) error {
	payload := LoginPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var token string
	var err error

	switch payload.Type {
	case GOOGLE:
		token, err = h.loginGoogle(payload.Token)
	case LINKEDIN:
		token, err = h.loginLinkedIn(payload.Token)
	case EMAIL:
		token, err = h.loginEmail(payload.Email, payload.Password)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		c.Logger().Error(err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "User or password is not correct")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
