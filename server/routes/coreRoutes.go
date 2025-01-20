package router

import (
	"github.com/Aman123at/cdc-go/controllers"
	"github.com/gin-gonic/gin"
)

func CdcRoutes(router *gin.Engine) {

	api := router.Group("/api")
	{
		api.POST("/create/table", controllers.CreateTable)
		api.GET("/fetch/tables", controllers.GetAllTablesData)
		api.GET("/fetch/collections", controllers.GetAllCollectionsData)
		api.POST("/delete/row", controllers.DeleteRowController)
		api.PUT("/edit/row", controllers.EditRowController)
		api.POST("/insert/row", controllers.InsertRowController)
	}
}
