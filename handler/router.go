package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(awshandler AwsProjHandler) (*Router, error) {

	router := gin.New()
	router.RedirectTrailingSlash = false

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
		})
	})

	router.Use(
		gin.Recovery(),
	)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/aws-api")
	{
		api.POST("")
		api.POST("/upload", awshandler.AwsProjUpload)
	}

	return &Router{router}, nil
}
