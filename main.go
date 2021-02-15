package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/Arun4rangan/api-tutorme/src/auth"
)

// Validator for echo
type Validator struct {
	validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func getPostgresURI() string {
	postgresURI := os.Getenv("POSTGRES_URI")
	strByte, err := ioutil.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))

	if err == nil {
		password := string(strByte)
		search := "__PASSWORD__"
		postgresURI = strings.Replace(postgresURI, search, password, 1)
	}
	return postgresURI
}

func main() {
	signingKey, err := auth.GetSigningKey()

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	_, err = auth.GetVerifyingKey()

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	// Validator
	validate := validator.New()

	validate.RegisterStructValidation(auth.LoginPayloadValidation, auth.LoginPayload{})
	validate.RegisterStructValidation(auth.SignUpPayloadValidation, auth.SignUpPayload{})

	db, err := sqlx.Connect("pgx", getPostgresURI())

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	e := echo.New()

	e.Use(middleware.Logger())

	// Prints stack trace and handles the
	// control to the centralized HTTPErrorHandler
	e.Use(middleware.Recover())

	// Body Limit Middleware
	e.Use(middleware.BodyLimit("10M"))

	e.Validator = &Validator{validator: validate}

	au := auth.NewStore(db)
	authHandler := auth.NewHandler(*au, signingKey)
	authHandler.Register(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	e.Logger.SetLevel(log.DEBUG)

	s := &http.Server{
		Addr:         string(":8010"),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	e.Logger.Debug(e.StartServer(s))
}
