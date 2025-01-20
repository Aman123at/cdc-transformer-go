package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // Specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Specify allowed headers
		ExposeHeaders:    []string{"Content-Length"},                          // Specify exposed headers
		AllowCredentials: true,                                                // Allow credentials (cookies, authorization headers)
	}))
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to CDC"})
	})
	CdcRoutes(router)
	return router
}
