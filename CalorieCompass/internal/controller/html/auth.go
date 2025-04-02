package html

import (
	
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"CalorieCompass/internal/entity"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, input entity.UserSignUp) (entity.AuthResponse, error)
	Login(ctx context.Context, input entity.UserLogin) (entity.AuthResponse, error)
}

type AuthController struct {
	authUseCase AuthUseCase
}

func NewAuthController(authUseCase AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

func (c *AuthController) ShowLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func (c *AuthController) ShowRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (c *AuthController) ProcessLogin(ctx *gin.Context) {
	var input entity.UserLogin
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login",
			"error": err.Error(),
		})
		return
	}

	response, err := c.authUseCase.Login(ctx.Request.Context(), input)
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	ctx.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) ProcessRegister(ctx *gin.Context) {
	var input entity.UserSignUp
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": err.Error(),
		})
		return
	}

	response, err := c.authUseCase.SignUp(ctx.Request.Context(), input)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register",
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	ctx.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
	ctx.Redirect(http.StatusSeeOther, "/login")
}