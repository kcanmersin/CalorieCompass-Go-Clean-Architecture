package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"CalorieCompass/internal/entity"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id int64) (entity.UserResponse, error)
}

type UserController struct {
	userUseCase UserUseCase
}

func NewUserController(userUseCase UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// @Summary Get user by ID
// @Description Get user information by ID
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := c.userUseCase.GetUserByID(ctx.Request.Context(), userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}