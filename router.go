package galleryapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Static("/media", "./media")

	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:    []string{"Authorization", "Content-Type"},
	}
	c := cors.New(config)
	r.Use(c)

	r.POST("/register", userRegisterHandler)
	r.POST("/get-token", userGetTokenHandler)

	userGroup := r.Group("/user")
	userGroup.Use(authentication())
	userController(userGroup)

	galleryGroup := r.Group("/galleries")
	galleryGroup.Use(authentication())
	galleryController(galleryGroup)

	return r
}
