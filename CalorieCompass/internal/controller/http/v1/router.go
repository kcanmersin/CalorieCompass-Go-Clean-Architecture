package v1

import (
	"CalorieCompass/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
)

type TokenValidator interface {
	ValidateToken(token string) (int64, error)
}

func NewRouter(handler *gin.Engine, authController *AuthController, userController *UserController, foodController *FoodController, tokenRepo TokenValidator) {
	// Create two route groups:
	// 1. Routes for the API with the /api/v1 prefix (for backwards compatibility)
	apiV1 := handler.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/sign-up", authController.SignUp)
			auth.POST("/login", authController.Login)
		}

		user := apiV1.Group("/user")
		user.Use(middleware.JWTAuth(tokenRepo))
		{
			user.GET("", userController.GetUserInfo)
		}

		// Food routes
		food := apiV1.Group("/food")
		{
			food.GET("/search", foodController.SearchFoods)
			food.GET("/:food_id", foodController.GetFoodDetails)
		}
	}

	// 2. Routes without the /api/v1 prefix (for Swagger to work correctly)
	auth := handler.Group("/auth")
	{
		auth.POST("/sign-up", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	user := handler.Group("/user")
	user.Use(middleware.JWTAuth(tokenRepo))
	{
		user.GET("", userController.GetUserInfo)
	}

	// Food routes
	food := handler.Group("/food")
	{
		food.GET("/search", foodController.SearchFoods)
		food.GET("/:food_id", foodController.GetFoodDetails)
	}
}
