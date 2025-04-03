// auth.go
package v1

import (
	"context"
	"net/http"

	"CalorieCompass/internal/entity"
	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, input entity.UserSignUp) (entity.AuthResponse, error)
	Login(ctx context.Context, input entity.UserLogin) (entity.AuthResponse, error)
	GetUserByID(ctx context.Context, id int64) (entity.UserResponse, error)
}

type AuthController struct {
	authUseCase AuthUseCase
}

func NewAuthController(authUseCase AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entity.UserSignUp true "User signup info"
// @Success 201 {object} entity.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/sign-up [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	var input entity.UserSignUp
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authUseCase.SignUp(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary Login user
// @Description Login a user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body entity.UserLogin true "User login info"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input entity.UserLogin
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authUseCase.Login(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
