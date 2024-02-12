package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterAuthRoutes register auth routes
func RegisterAuthRoutes(e *echo.Echo, validate *validator.Validate, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, authUseCase rfrl.AuthUseCase) {
	authview := views.AuthView{
		AuthUseCases: authUseCase,
		Key:          *privateKey,
	}
	validate.RegisterStructValidation(views.LoginPayloadValidation, views.LoginPayload{})
	e.POST("/login/", authview.Login)

	validate.RegisterStructValidation(views.SignUpPayloadValidation, views.SignUpPayload{})
	e.POST("/signup/", authview.Signup)

	loginAuth := e.Group("/login-authorized")
	loginAuth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	loginAuth.POST("/", authview.AuthorizedLogin)

	signUpFlowView := e.Group("/sign-up-flow")
	signUpFlowView.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	signUpFlowView.PUT("/", authview.UpdateSignUpFlow)

	blockAuth := e.Group("/block")
	blockAuth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	blockAuth.POST("/", authview.BlockClient)

}
