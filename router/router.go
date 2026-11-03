package router

import (
	"app/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	object := api.Group("/object")
	{
		object.GET("/list", handler.ListObjectsHandler)
		object.GET("/info", handler.GetObjectInfoHandler)
		object.DELETE("/delete", handler.DeleteObjectHandler)
		object.POST("/upload", handler.UploadObjectHandler)
	}

	return router
}
