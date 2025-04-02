package html

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"CalorieCompass/internal/entity"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id int64) (entity.UserResponse, error)
}

type TokenValidator interface {
	ValidateToken(token string) (int64, error)
}

type HomeController struct {
	userUseCase  UserUseCase
	tokenService TokenValidator
}

func NewHomeController(userUseCase UserUseCase, tokenService TokenValidator) *HomeController {
	return &HomeController{
		userUseCase:  userUseCase,
		tokenService: tokenService,
	}
}

func (c *HomeController) Home(ctx *gin.Context) {
	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
			"user":  nil,
		})
		return
	}

	userID, err := c.tokenService.ValidateToken(token)
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
			"user":  nil,
		})
		return
	}

	user, err := c.userUseCase.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
			"user":  nil,
		})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home",
		"user":  user,
	})
}