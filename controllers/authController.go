package controllers

import (
	"fmt"
	"joranvest/dto"
	"joranvest/helper"
	entity_view_models "joranvest/models/entity_view_models"
	"joranvest/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	Register(ctx *gin.Context)
	RegisterForm(ctx *gin.Context)
}

type authController struct {
	//-- Here to put your service
	authService  service.AuthService
	emailService service.EmailService
	jwtService   service.JWTService
	ginService   *gin.Context
}

var (
	ctx *gin.Context
	//userSession helper.UserSession = helper.NewUserSession(ctx)
)

//-- to create a new instance of AuthController
func NewAuthController(authService service.AuthService, emailService service.EmailService, jwtService service.JWTService) AuthController {
	return &authController{
		authService:  authService,
		emailService: emailService,
		jwtService:   jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDto dto.LoginDto
	err := ctx.ShouldBind(&loginDto)
	fmt.Println(loginDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request login", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDto.Username, loginDto.Email, loginDto.Password)

	if !authResult.Status {
		response := helper.BuildResponse(false, authResult.Message, helper.EmptyObj{})
		ctx.JSON(http.StatusOK, response)
		return
	}

	if v, ok := (authResult.Data).(entity_view_models.EntityApplicationUserView); ok {
		// generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		generatedToken := c.jwtService.GenerateToken(v.Id, v.EntityId)
		v.Token = generatedToken

		response := helper.BuildResponse(true, "Ok!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildErrorResponse("Please check your credentials", "Invalid Credentials", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Logout(ctx *gin.Context) {
	println("==========================")
	println("==========================")
	println("==========================")
	println("Logout...")
	session := sessions.Default(ctx)
	session.Clear() // this will mark the session as "written" only if there's
	// at least one key to delete
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	ctx.Redirect(301, "/")
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDto dto.ApplicationUserRegisterDto
	err := ctx.ShouldBind(&registerDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request register", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDto.Email) {
		response := helper.BuildResponse(false, "Email telah terdaftar.", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, response)
		return
	} else {
		createdUser, err := c.authService.CreateUser(registerDto)
		if err != nil {
			var message = ""
			if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "idx_users_username") {
				message = "Username " + registerDto.Email + " sudah terdaftar"
				response := helper.BuildErrorResponse("Failed to register user", message, helper.EmptyObj{})
				ctx.JSON(http.StatusConflict, response)
			} else {
				message = fmt.Sprintf("%v", err.Error())
				response := helper.BuildErrorResponse("Failed to register user", message, helper.EmptyObj{})
				ctx.JSON(http.StatusBadRequest, response)
			}
		} else {

			token := c.jwtService.GenerateToken(createdUser.Id, createdUser.EntityId)

			// Send Email Verification
			to := []string{createdUser.Email}
			c.emailService.SendEmailVerification(to, createdUser.Id)

			createdUser.Token = token
			response := helper.BuildResponse(true, "Ok!", createdUser)
			ctx.JSON(http.StatusOK, response)
		}
		return
	}
}

func (c *authController) RegisterForm(ctx *gin.Context) {
	var registerDto dto.ApplicationUserRegisterDto
	err := ctx.ShouldBind(&registerDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request register", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDto.Email) {
		response := helper.BuildErrorResponse("Failed to request register", "Duplicate Email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser, err := c.authService.CreateUser(registerDto)
		if err != nil {
			var message = ""
			if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "idx_users_username") {
				message = "Username " + registerDto.Email + " sudah terdaftar"
				response := helper.BuildErrorResponse("Failed to register user", message, helper.EmptyObj{})
				ctx.JSON(http.StatusConflict, response)
			} else {
				message = fmt.Sprintf("%v", err.Error())
				response := helper.BuildErrorResponse("Failed to register user", message, helper.EmptyObj{})
				ctx.JSON(http.StatusBadRequest, response)
			}
		} else {

			token := c.jwtService.GenerateToken(createdUser.Id, createdUser.EntityId)

			createdUser.Token = token
			response := helper.BuildResponse(true, "Ok!", createdUser)
			ctx.JSON(http.StatusOK, response)
		}
		return
	}
}
