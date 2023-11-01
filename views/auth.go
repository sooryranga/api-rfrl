package views

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type loginFields struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Token    string `json:"token"`
	Password string `json:"password" validate:"omitempty,gte=10"`
	Type     string `json:"type" validate:"required,oneof= google linkedin email"`
}

type (
	// SignUpPayload is the struct used to hold payload from /signup
	SignUpPayload struct {
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Photo     string    `json:"profileImageURL"`
		About     string    `json:"about"`
		IsTutor   null.Bool `json:"isTutor"`
		loginFields
	}
	// LoginPayload is the struct used to hold payload from /login
	LoginPayload struct {
		loginFields
	}
)

// SignUpPayloadValidation validates client inputs
func SignUpPayloadValidation(sl validator.StructLevel) {
	payload := sl.Current().Interface().(SignUpPayload)

	switch payload.Type {
	case tutorme.GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case tutorme.LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case tutorme.EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// LoginPayloadValidation validates client inputs
func LoginPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(LoginPayload)

	switch payload.Type {
	case tutorme.GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case tutorme.LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case tutorme.EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

type AuthView struct {
	AuthUseCases tutorme.AuthUseCase
	Key          rsa.PrivateKey
}

// Signup endpoint
func (av *AuthView) Signup(c echo.Context) error {
	payload := SignUpPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var newClient *tutorme.Client
	var err error

	switch payload.Type {
	case tutorme.GOOGLE:
		newClient, err = av.AuthUseCases.SignupGoogle(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case tutorme.LINKEDIN:
		newClient, err = av.AuthUseCases.SignupLinkedIn(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case tutorme.EMAIL:
		newClient, err = av.AuthUseCases.SignupEmail(
			payload.Password,
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		c.Logger().Error(err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return av.login(c, payload.loginFields)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}

	claims := &tutorme.JWTClaims{
		newClient.ID,
		newClient.Email.String,
		newClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(claims, &av.Key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token,
		"client": newClient,
	})
}

func (av *AuthView) AuthorizedLogin(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingClient, err := av.AuthUseCases.LoginWithJWT(claims.ClientID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	newClaims := &tutorme.JWTClaims{
		existingClient.ID,
		existingClient.Email.String,
		existingClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(newClaims, &av.Key)

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token,
		"client": existingClient,
	})
}

func (av *AuthView) login(c echo.Context, payload loginFields) error {
	var existingClient *tutorme.Client
	var err error

	switch payload.Type {
	case tutorme.GOOGLE:
		existingClient, err = av.AuthUseCases.LoginGoogle(payload.Token)
	case tutorme.LINKEDIN:
		existingClient, err = av.AuthUseCases.LoginLinkedIn(payload.Token)
	case tutorme.EMAIL:
		existingClient, err = av.AuthUseCases.LoginEmail(payload.Email, payload.Password)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	claims := &tutorme.JWTClaims{
		existingClient.ID,
		existingClient.Email.String,
		existingClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(claims, &av.Key)

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token,
		"client": existingClient,
	})
}

// Login is used to login an client
func (av *AuthView) Login(c echo.Context) error {
	payload := LoginPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return av.login(c, payload.loginFields)
}
