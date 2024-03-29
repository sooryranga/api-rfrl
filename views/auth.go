package views

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
		Stage rfrl.SignUpFlow `json:"stage"`
	}

	BlockClientPayload struct {
		ClientID null.String `json:"clientId" validate:"required"`
		Blocked  null.Bool   `json:"blocked" validate:"required"`
	}
)

// SignUpPayloadValidation validates client inputs
func SignUpPayloadValidation(sl validator.StructLevel) {
	payload := sl.Current().Interface().(SignUpPayload)

	switch payload.Type {
	case rfrl.GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case rfrl.LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case rfrl.EMAIL:
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

	switch payload.Type {
	case rfrl.GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case rfrl.LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case rfrl.EMAIL:
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
	AuthUseCases rfrl.AuthUseCase
	Key          rsa.PrivateKey
}

// Signup endpoint
func (av *AuthView) Signup(c echo.Context) error {
	payload := SignUpPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "Signup - Bind"))
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(errors.Wrap(err, "Signup - Validate"))
	}
	var newClient *rfrl.Client
	var auth *rfrl.Auth
	var err error

	switch payload.Type {
	case rfrl.GOOGLE:
		newClient, auth, err = av.AuthUseCases.SignupGoogle(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case rfrl.LINKEDIN:
		newClient, auth, err = av.AuthUseCases.SignupLinkedIn(
			payload.Token,
			payload.Email,
			payload.FirstName,
			payload.LastName,
			payload.Photo,
			payload.About,
			payload.IsTutor,
		)
	case rfrl.EMAIL:
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return av.login(c, payload.loginFields)
			}
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	claims := &rfrl.JWTClaims{
		newClient.ID,
		newClient.Email.String,
		newClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(claims, &av.Key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token,
		"client": newClient,
		"auth":   auth,
	})
}

func (av *AuthView) AuthorizedLogin(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	existingClient, auth, err := av.AuthUseCases.LoginWithJWT(claims.ClientID)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, "Client not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
		}
	}

	newClaims := &rfrl.JWTClaims{
		existingClient.ID,
		existingClient.Email.String,
		existingClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(newClaims, &av.Key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token,
		"client": existingClient,
		"auth":   auth,
	})
}

func (av *AuthView) login(c echo.Context, payload loginFields) error {
	var existingClient *rfrl.Client
	var auth *rfrl.Auth
	var err error

	switch payload.Type {
	case rfrl.GOOGLE:
		existingClient, auth, err = av.AuthUseCases.LoginGoogle(payload.Token)
	case rfrl.LINKEDIN:
		existingClient, auth, err = av.AuthUseCases.LoginLinkedIn(payload.Token)
	case rfrl.EMAIL:
		existingClient, auth, err = av.AuthUseCases.LoginEmail(payload.Email, payload.Password)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, "Client doesn't exist in our records").SetInternal(err)
		case bcrypt.ErrMismatchedHashAndPassword:
			return echo.NewHTTPError(http.StatusNotFound, "Email and password do not match").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
		}
	}

	claims := &rfrl.JWTClaims{
		existingClient.ID,
		existingClient.Email.String,
		existingClient.IsAdmin.Bool,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token, err := av.AuthUseCases.GenerateToken(claims, &av.Key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return av.login(c, payload.loginFields)
}

func (av *AuthView) UpdateSignUpFlow(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	payload := SignUpFlowPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	err = av.AuthUseCases.UpdateSignUpFlow(claims.ClientID, payload.Stage)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (av *AuthView) BlockClient(c echo.Context) error {
	claims, err := rfrl.GetClaims(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if !claims.Admin {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are unauthorized to use this view")
	}

	payload := BlockClientPayload{}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	err = av.AuthUseCases.BlockClient(
		payload.ClientID.String,
		payload.Blocked.Bool,
	)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.NewHTTPError(http.StatusNotFound, "Client not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	}

	return c.NoContent(http.StatusOK)
}
