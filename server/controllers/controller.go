package controllers

import (
	"strings"

	"github.com/Aman123at/cdc-go/connections"
	"github.com/Aman123at/cdc-go/models"
	"github.com/Aman123at/cdc-go/utils"
	"github.com/gin-gonic/gin"
)

func CreateTable(c *gin.Context) {
	var body models.CreateTableReq

	if binderr := c.BindJSON(&body); binderr != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body.TableName = strings.TrimSpace(body.TableName)

	if body.TableName == "" {
		c.JSON(400, gin.H{"error": "Table name is required"})
		return
	}

	if len(body.TableName) > 25 {
		c.JSON(400, gin.H{"error": "Table name should be under 25 characters"})
		return
	}

	if len(body.Columns) == 0 {
		c.JSON(400, gin.H{"error": "Atleast one column is required"})
		return
	}

	for _, item := range body.Columns {
		utils.TrimColumnSpaces(item)
	}

	duplicationErr := body.ValidateNoDuplicateColumns()

	if duplicationErr != nil {
		c.JSON(400, gin.H{"error": duplicationErr})
		return
	}

	err := connections.CreateNewTable(body)

	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"message": "Created table successfully."})
}

type TableResponse struct {
	Name string `json:"name"`
	Rows []any  `json:"rows"`
}

func InsertRowController(c *gin.Context) {
	var body models.InsertRowReq

	if binderr := c.BindJSON(&body); binderr != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body.TableName = strings.TrimSpace(body.TableName)

	if body.TableName == "" {
		c.JSON(400, gin.H{"error": "Table name is required"})
		return
	}

	if len(body.Row) == 0 {
		c.JSON(400, gin.H{"error": "Row data is required"})
		return
	}

	err := connections.InsertRow(body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Row inserted successfully"})
}

func EditRowController(c *gin.Context) {
	var body models.EditRowReq

	if binderr := c.BindJSON(&body); binderr != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body.TableName = strings.TrimSpace(body.TableName)

	if body.TableName == "" {
		c.JSON(400, gin.H{"error": "Table name is required"})
		return
	}

	if body.RowId == 0 {
		c.JSON(400, gin.H{"error": "Row ID is required"})
		return
	}

	if len(body.Row) == 0 {
		c.JSON(400, gin.H{"error": "Update data is required"})
		return
	}

	whereCondition := map[string]any{"id": body.RowId}

	err := connections.EditRow(body, whereCondition)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Row updated successfully"})
}

func DeleteRowController(c *gin.Context) {
	var body models.DeleteRowReq

	if binderr := c.BindJSON(&body); binderr != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body.TableName = strings.TrimSpace(body.TableName)

	if body.TableName == "" {
		c.JSON(400, gin.H{"error": "Table name is required"})
		return
	}

	if body.RowId == 0 {
		c.JSON(400, gin.H{"error": "Row ID is required"})
		return
	}

	whereCondition := map[string]any{"id": body.RowId}

	err := connections.DeleteRow(body.TableName, whereCondition)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Row deleted successfully"})
}

func GetAllTablesData(c *gin.Context) {
	data, err := connections.GetAllTablesData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}

func GetAllCollectionsData(c *gin.Context) {
	data, err := connections.GetAllCollectionsData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
