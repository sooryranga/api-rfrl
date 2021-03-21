package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTutorReviewRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, tutorReviewUseCase tutorme.TutorReviewUseCase) {
	tutorReviewView := views.TutorReviewView{TutorReviewUseCase: tutorReviewUseCase}

	tutorReviewR := e.Group("/tutor-review")
	tutorReviewR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	tutorReviewR.POST("/", tutorReviewView.CreateTutorReviewEndpoint)
	tutorReviewR.PUT("/:id/", tutorReviewView.UpdateTutorReviewEndpoint)
	tutorReviewR.DELETE("/:id/", tutorReviewView.DeleteTutorReviewEndpoint)
	tutorReviewR.GET("/:id/", tutorReviewView.GetTutorReviewEndpoint)

	tutorReviewsR := e.Group("/tutor-reviews")
	tutorReviewsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	tutorReviewsR.GET("/:tutorID/", tutorReviewView.GetTutorReviewsEndpoint)

	tutorReviewsAggregateR := e.Group("/tutor-reviews-aggregate")
	tutorReviewsAggregateR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	tutorReviewsAggregateR.GET("/:tutorID/", tutorReviewView.GetTutorReviewsAggregateEndpoint)
}
