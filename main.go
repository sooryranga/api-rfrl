package main

import (
	"fmt"
	"net/http"
	"time"

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

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func main() {
	//Config
	var cfg Config
	readEnv(&cfg)
	fmt.Printf("%+v", cfg)

	e := echo.New()

	e.Use(middleware.Logger())

	// Prints stack trace and handles the
	// control to the centralized HTTPErrorHandler
	e.Use(middleware.Recover())

	// Body Limit Middleware
	e.Use(middleware.BodyLimit("10M"))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/login", auth.Login)
	e.Logger.SetLevel(log.DEBUG)

	s := &http.Server{
		Addr:         ":8010",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	e.Logger.Fatal(e.StartServer(s))
}
