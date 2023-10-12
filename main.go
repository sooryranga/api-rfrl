package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Arun4rangan/api-tutorme/routes"
	"github.com/Arun4rangan/api-tutorme/store"
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/usecases"
	"github.com/go-playground/validator"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
	signingKey, err := tutorme.GetSigningKey()

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	publicKey, err := tutorme.GetVerifyingKey()

	fmt.Println(fmt.Sprintf("%v", publicKey))

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	// Validator
	validate := validator.New()

	db, err := sqlx.Connect("pgx", getPostgresURI())
	defer db.Close()

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

	// Stores
	authStore := store.NewAuthStore()
	clientStore := store.NewClientStore()
	documentStore := store.NewDocumentStore()
	sessionStore := store.NewSessionStore()
	tutorReviewStore := store.NewTutorReviewStore()

	// Usecases
	authUseCase := usecases.NewAuthUseCase(*db, authStore, clientStore)
	clientUseCase := usecases.NewClientUseCase(*db, clientStore)
	documentUseCase := usecases.NewDocumentUseCase(*db, documentStore)
	sessionUseCase := usecases.NewSessionUseCase(*db, sessionStore)
	tutorUseCase := usecases.NewTutorReviewUseCase(db, tutorReviewStore, sessionStore, clientStore)

	routes.RegisterAuthRoutes(e, validate, signingKey, authUseCase)
	routes.RegisterClientRoutes(e, validate, publicKey, clientUseCase)
	routes.RegisterDocumentRoutes(e, validate, publicKey, documentUseCase)
	routes.RegisterSessionRoutes(e, validate, publicKey, sessionUseCase)
	routes.RegisterTutorReviewRoutes(e, validate, publicKey, tutorUseCase)

	e.Validator = &Validator{validator: validate}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	e.Logger.SetLevel(log.DEBUG)

	s := &http.Server{
		Addr:         string(":8010"),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	e.Logger.Fatal(e.StartServer(s))
}
