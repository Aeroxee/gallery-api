package galleryapi

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Aeroxee/gallery-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func galleryGetHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"message":    "Authentication is required.",
			"error_type": "authentication_error",
		})
		return
	}

	limit := getQueryInt(ctx.Request, "limit", 10)
	offset := getQueryInt(ctx.Request, "offset", 1)
	offsetQ := (offset - 1) * limit

	galleries := models.GetAllGallery(user.ID, limit, offsetQ)
	ctx.JSON(http.StatusOK, galleries)
}

func galleryDetailHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"error_type": "authentication_error",
			"message":    "Authentication is required.",
		})
		return
	}

	galleryId := ctx.Param("galleryId")
	galleryIdInt, err := strconv.Atoi(galleryId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "Please set gallery id with integer type.",
			"error_type": "type_error",
		})
		return
	}

	gallery, err := models.GetGalleryByID(galleryIdInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	// check user is owner
	if user.ID != gallery.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "You don't have permission to read this photo.",
			"error_type": "permission_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gallery)
}

func galleryCreateHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"error_type": "authentication_error",
			"message":    "Authentication is required.",
		})
		return
	}

	payloads := struct {
		Photo       *multipart.FileHeader `form:"photo" validate:"required"`
		Title       string                `form:"title" validate:"required"`
		Description string                `form:"description"`
	}{}
	err = ctx.ShouldBindWith(&payloads, binding.FormMultipart)
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
			fmt.Println(err.Error())
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

	// check extension file
	ext := filepath.Ext(payloads.Photo.Filename)
	if !isAllowedExtension(ext) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"error_type": "extenstion_error",
			"message":    "Please upload a file with jpeg|jpg|png|webp extension only.",
		})
		return
	}

	gallery := models.Gallery{
		Title:       payloads.Title,
		Description: &payloads.Description,
		UserID:      user.ID,
	}

	err = models.CreateNewGallery(&gallery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"error_type": "create_new_gallery_error",
			"message":    err.Error(),
		})
		return
	}

	filename := payloads.Photo.Filename
	destination := fmt.Sprintf("media/galleries/%s/%d/%s", user.Username, gallery.ID, filename)
	// upload
	err = ctx.SaveUploadedFile(payloads.Photo, destination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"message":    err.Error(),
			"error_type": "upload_error",
		})
		return
	}

	gallery.Photo = &destination
	models.DB().Save(&gallery)

	ctx.JSON(http.StatusCreated, gallery)
}

func galleryUpdateHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"error_type": "authentication_error",
			"message":    "Authentication is required.",
		})
		return
	}

	galleryId := ctx.Param("galleryId")
	galleryIdInt, err := strconv.Atoi(galleryId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "Please set gallery id with integer type.",
			"error_type": "type_error",
		})
		return
	}

	gallery, err := models.GetGalleryByID(galleryIdInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	// check user is owner
	if user.ID != gallery.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "You don't have permission to update this photo.",
			"error_type": "permission_error",
		})
		return
	}

	payloads := struct {
		Title       string                `form:"title"`
		Description string                `form:"description"`
		Photo       *multipart.FileHeader `form:"photo"`
	}{}
	err = ctx.ShouldBindWith(&payloads, binding.FormMultipart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    err.Error(),
			"error_type": "payload_error",
		})
		return
	}

	if payloads.Title != "" {
		gallery.Title = payloads.Title
	}
	if payloads.Description != "" {
		gallery.Description = &payloads.Description
	}
	if payloads.Photo != nil {
		if !isAllowedExtension(filepath.Ext(payloads.Photo.Filename)) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":     "error",
				"error_type": "extenstion_error",
				"message":    "Please upload a file with jpeg|jpg|png|webp extension only.",
			})
			return
		}

		oldFile := gallery.Photo
		// remove old file
		err = os.RemoveAll(*oldFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":     "error",
				"message":    err.Error(),
				"error_type": "server_error",
			})
			return
		}

		filename := payloads.Photo.Filename
		newFile := fmt.Sprintf("media/galleries/%s/%d/%s", user.Username, gallery.ID, filename)

		err = ctx.SaveUploadedFile(payloads.Photo, newFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":     "error",
				"message":    err.Error(),
				"error_type": "server_error",
			})
			return
		}

		gallery.Photo = &newFile
	}

	models.DB().Save(&gallery)
	ctx.JSON(http.StatusOK, gallery)
}

func galleryDeleteHandler(ctx *gin.Context) {
	user, err := getUserContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":     "error",
			"error_type": "authentication_error",
			"message":    "Authentication is required.",
		})
		return
	}

	galleryId := ctx.Param("galleryId")
	galleryIdInt, err := strconv.Atoi(galleryId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "Please set gallery id with integer type.",
			"error_type": "type_error",
		})
		return
	}

	gallery, err := models.GetGalleryByID(galleryIdInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	// check user is owner
	if user.ID != gallery.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"message":    "You don't have permission to delete this photo.",
			"error_type": "permission_error",
		})
		return
	}

	models.DB().Delete(&gallery)
	ctx.JSON(http.StatusNoContent, nil)
}
