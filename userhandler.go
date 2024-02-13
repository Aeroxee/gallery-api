package galleryapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Aeroxee/gallery-api/auth"
	"github.com/Aeroxee/gallery-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func userRegisterHandler(ctx *gin.Context) {
	payloads := struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"error_type": "payload_error",
			"message":    err.Error(),
		})
		return
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&payloads)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Error on field: %s, with %s.", err.Field(), err.ActualTag()))
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"error_type": "validation_error",
			"messages":   errorMessages,
		})
		return
	}

	user := models.User{
		FirstName: payloads.FirstName,
		LastName:  payloads.LastName,
		Username:  payloads.Username,
		Email:     payloads.Email,
		Password:  payloads.Password,
	}

	// save
	err = models.CreateNewUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    err.Error(),
			"error_type": "register_error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func userGetTokenHandler(ctx *gin.Context) {
	payloads := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"error_type": "payload_error",
			"message":    err.Error(),
		})
		return
	}

	var user models.User
	if strings.Contains(payloads.Username, "@") {
		user, err = models.GetUserByEmail(payloads.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":     "error",
				"message":    "Username or password is incorrect.",
				"error_type": "gettoken_error",
			})
			return
		}
	} else {
		user, err = models.GetUserByUsername(payloads.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":     "error",
				"message":    "Username or password is incorrect.",
				"error_type": "gettoken_error",
			})
			return
		}
	}

	credential := auth.Credential{
		UserID: user.ID,
	}
	token, err := auth.GetToken(credential)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"error_type": "gettoken_error",
			"message":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Generate new token is successfully.",
		"token":      token,
		"error_type": nil,
	})
}

func userAuthHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"message":    "Authentication is required.",
			"error_type": "authorize_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
