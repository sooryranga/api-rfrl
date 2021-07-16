package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTutorReviewRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, tutorReviewUseCase rfrl.TutorReviewUseCase) {
	tutorReviewView := views.TutorReviewView{TutorReviewUseCase: tutorReviewUseCase}

	tutorReviewR := e.Group("/tutor-review")
	tutorReviewR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	tutorReviewR.POST("/", tutorReviewView.CreateTutorReviewEndpoint)
	tutorReviewR.PUT("/:id/", tutorReviewView.UpdateTutorReviewEndpoint)
	tutorReviewR.DELETE("/:id/", tutorReviewView.DeleteTutorReviewEndpoint)
	tutorReviewR.GET("/:id/", tutorReviewView.GetTutorReviewEndpoint)

	tutorReviewsR := e.Group("/tutor-reviews")
	tutorReviewsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	tutorReviewsR.GET("/:tutorID/", tutorReviewView.GetTutorReviewsEndpoint)

	tutorReviewsAggregateR := e.Group("/tutor-reviews-aggregate")
	tutorReviewsAggregateR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	tutorReviewsAggregateR.GET("/:tutorID/", tutorReviewView.GetTutorReviewsAggregateEndpoint)

	pendingReviewsR := e.Group("/pending-tutor-reviews")
	pendingReviewsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	pendingReviewsR.GET("/", tutorReviewView.GetPendingReviewsEndpoint)

}
