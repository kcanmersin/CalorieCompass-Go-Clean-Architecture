package v1

import (
	"github.com/gin-gonic/gin"
	"CalorieCompass/internal/controller/http/middleware"
)

type TokenValidator interface {
	ValidateToken(token string) (int64, error)
}

func NewRouter(handler *gin.Engine, authController *AuthController, tokenRepo TokenValidator) {
	h := handler.Group("/api/v1")

	auth := h.Group("/auth")
	{
		auth.POST("/sign-up", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	user := h.Group("/user")
	user.Use(middleware.JWTAuth(tokenRepo))
	{
		// Protected routes
	}
}