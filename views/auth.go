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
)

type (
	// SignUpPayload is the struct used to hold payload from /signup
	SignUpPayload struct {
		Email     string `json:"email" validate:"omitempty,email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Token     string `json:"token"`
		Photo     string `json:"profileImageURL"`
		About     string `json:"about"`
		Password  string `json:"password" validate:"omitempty,gte=10"`
		Type      string `json:"type" validate:"required,oneof= google linkedin email"`
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

type (
	// LoginPayload is the struct used to hold payload from /login
	LoginPayload struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Token    string `json:"token"`
		Password string `json:"password" validate:"omitempty,gte=10"`
		Type     string `json:"type" validate:"required,oneof= google linkedin email"`
	}
)

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
		)
	case tutorme.LINKEDIN:
		newClient, err = av.AuthUseCases.SignupLinkedIn(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
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
		)
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
					fmt.Sprintf("Client is already signed up with %s", payload.Type),
				)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong")
	}

	claims := &tutorme.JWTClaims{
		newClient.ID,
		newClient.Email.String,
		false,
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
		false,
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

// Login is used to login an client
func (av *AuthView) Login(c echo.Context) error {
	payload := LoginPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
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
		c.Logger().Error(err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Client or password is not correct")
	}

	claims := &tutorme.JWTClaims{
		existingClient.ID,
		existingClient.Email.String,
		false,
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
