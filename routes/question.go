package routes

import (
	"crypto/rsa"

	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisteerQuestionRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, questionUseCase tutorme.QuestionUseCase) {
	questionViews := views.QuestionView{QuestionUseCase: questionUseCase}

	questionR := e.Group("/question")
	questionR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	questionR.POST("/", questionViews.CreateQuestionEndpoint)
	questionR.GET("/:id/", questionViews.GetQuestionEndpoint)
	questionR.DELETE("/:id/", questionViews.DeleteQuestionEndpoint)
	questionR.PUT("/:id/", questionViews.UpdateQuestionEndpoint)

	questionsR := e.Group("/questions")
	questionsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	questionsR.GET("/", questionViews.GetQuestionsEndpoint)
	questionsR.GET("/:id/", questionViews.GetQuestionsFromClientEndpoint)

	applyToQuestionR := e.Group("/question/:questionID/apply")
	applyToQuestionR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))
	applyToQuestionR.POST("/", questionViews.ApplyToQuestionEndpoint)
}
