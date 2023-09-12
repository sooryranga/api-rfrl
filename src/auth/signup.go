package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	// SignUpPayload is the struct used to hold payload from /signup
	SignUpPayload struct {
		Email           string `json:"email" validate:"omitempty,email"`
		Name            string `json:"name"`
		Token           string `json:"token"`
		ProfileImageURL string `json:"profileImageURL"`
		Password        string `json:"password" validate:"omitempty,gte=10"`
		Type            string `json:"type" validate:"required,oneof= google linkedin email"`
	}
)

func hashAndSalt(pwd []byte) ([]byte, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return hash, nil
}

func (h *Handler) signupGoogle(payload SignUpPayload) (string, error) {
	_, err := h.authStore.CreateWithToken(payload.Token, GOOGLE)

	if err != nil {
		return "", err
	}

	claims := &jwtCustomClaims{
		payload.Email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return GenerateToken(claims)
}

func (h *Handler) signupLinkedIn(payload SignUpPayload) (string, error) {
	_, err := h.authStore.CreateWithToken(payload.Token, LINKEDIN)

	if err != nil {
		return "", err
	}

	claims := &jwtCustomClaims{
		"",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return GenerateToken(claims)
}

func (h *Handler) signupEmail(payload SignUpPayload) (string, error) {
	email := payload.Email

	hash, err := hashAndSalt([]byte(payload.Password))

	if err != nil {
		return "", err
	}

	_, err = h.authStore.CreateWithEmail(email, hash)

	if err != nil {
		return "", err
	}

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	return GenerateToken(claims)
}

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
		token, err = h.signupGoogle(payload)
	case LINKEDIN:
		token, err = h.signupLinkedIn(payload)
	case EMAIL:
		token, err = h.signupEmail(payload)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Login type - %s is not supported", payload.Token))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ErrNoRows")
		}
		var pgErr *pgconn.PgError
		fmt.Printf("%v", errors.As(err, &pgErr))
		if errors.As(err, &pgErr) {
			c.Logger().Error(pgErr)
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
