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
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type loginFields struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Token    string `json:"token"`
	Password string `json:"password" validate:"omitempty,gte=6"`
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

	SignUpFlowPayload struct {
		Stage tutorme.SignUpFlow `json:"stage"`
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
		if len(payload.Password) < 6 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// LoginPayloadValidation validates client inputs
func LoginPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(LoginPayload)
	log.Errorj(log.JSON{"payload": payload})
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
		if len(payload.Password) < 6 {
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
	var auth *tutorme.Auth
	var err error

	switch payload.Type {
	case tutorme.GOOGLE:
		newClient, auth, err = av.AuthUseCases.SignupGoogle(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case tutorme.LINKEDIN:
		newClient, auth, err = av.AuthUseCases.SignupLinkedIn(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case tutorme.EMAIL:
		newClient, auth, err = av.AuthUseCases.SignupEmail(
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
		"auth":   auth,
	})
}

func (av *AuthView) AuthorizedLogin(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingClient, auth, err := av.AuthUseCases.LoginWithJWT(claims.ClientID)

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
		"auth":   auth,
	})
}

func (av *AuthView) login(c echo.Context, payload loginFields) error {
	var existingClient *tutorme.Client
	var auth *tutorme.Auth
	var err error

	switch payload.Type {
	case tutorme.GOOGLE:
		existingClient, auth, err = av.AuthUseCases.LoginGoogle(payload.Token)
	case tutorme.LINKEDIN:
		existingClient, auth, err = av.AuthUseCases.LoginLinkedIn(payload.Token)
	case tutorme.EMAIL:
		existingClient, auth, err = av.AuthUseCases.LoginEmail(payload.Email, payload.Password)
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
		"auth":   auth,
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

func (av *AuthView) UpdateSignUpFlow(c echo.Context) error {
	claims, err := tutorme.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	payload := SignUpFlowPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return av.AuthUseCases.UpdateSignUpFlow(claims.ClientID, payload.Stage)
}
