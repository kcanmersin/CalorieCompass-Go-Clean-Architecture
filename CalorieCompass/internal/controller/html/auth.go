package html

import (
	"context"
	"fmt"
	"log"
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
	// Log form data to debug
	err := ctx.Request.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	log.Printf("Login Form data: %v", ctx.Request.PostForm)

	// Explicitly extract form fields
	email := ctx.PostForm("Email")
	password := ctx.PostForm("Password")
	
	if email == "" || password == "" {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login",
			"error": "Email and password are required",
		})
		return
	}

	input := entity.UserLogin{
		Email: email,
		Password: password,
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
	// Log form data to debug
	err := ctx.Request.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	log.Printf("Register Form data: %v", ctx.Request.PostForm)

	// Explicitly extract form fields
	name := ctx.PostForm("Name")
	email := ctx.PostForm("Email")
	password := ctx.PostForm("Password")
	
	if name == "" || email == "" || password == "" {
		errorMsg := fmt.Sprintf("Missing required fields. Name: %s, Email: %s, Password: [hidden]", name, email)
		ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": errorMsg,
		})
		return
	}

	input := entity.UserSignUp{
		Name: name,
		Email: email,
		Password: password,
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