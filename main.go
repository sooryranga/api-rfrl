package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/Arun4rangan/api-tutorme/src/auth"
)

// Server configs from ENV variables
type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
		Host string `envconfig:"SERVER_HOST"`
	}
	Database struct {
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
	}
}

// Validator for echo
type Validator struct {
	validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func main() {
	port := flag.String("port", ":8010", "Server listener port")
	serverName := flag.String("name", "Server1", "Server Name")
	flag.Parse()

	fmt.Printf("my port: \"%v\"\n", string(*port))

	//Config
	var cfg Config
	readEnv(&cfg)
	fmt.Printf("CONFIG : %+v", cfg)

	// Validator
	validate := validator.New()

	validate.RegisterStructValidation(auth.LoginPayloadValidation, auth.LoginPayload{})
	validate.RegisterStructValidation(auth.SignUpPayloadValidation, auth.SignUpPayload{})

	e := echo.New()

	e.Use(middleware.Logger())

	// Prints stack trace and handles the
	// control to the centralized HTTPErrorHandler
	e.Use(middleware.Recover())

	// Body Limit Middleware
	e.Use(middleware.BodyLimit("10M"))

	e.Validator = &Validator{validator: validate}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Hello from %v", string(*serverName)))
	})

	e.POST("/login", auth.Login)
	e.Logger.SetLevel(log.DEBUG)

	s := &http.Server{
		Addr:         string(*port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	e.Logger.Debug(e.StartServer(s))
}
