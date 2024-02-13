package galleryapi

import "github.com/gin-gonic/gin"

func userController(group *gin.RouterGroup) {
	group.GET("/auth", userAuthHandler)
}

func galleryController(group *gin.RouterGroup) {
	group.GET("", galleryGetHandler)
	group.POST("", galleryCreateHandler)
	group.GET("/:galleryId", galleryDetailHandler)
	group.PUT("/:galleryId/update", galleryUpdateHandler)
	group.DELETE("/:galleryId/delete", galleryDeleteHandler)
}
