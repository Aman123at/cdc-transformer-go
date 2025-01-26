package models

import (
	"fmt"
	"strings"
)

// DataType represents the supported SQL column data types
type DataType string

const (
	TypeInt       DataType = "int"
	TypeText      DataType = "text"
	TypeBoolean   DataType = "boolean"
	TypeVarchar   DataType = "varchar(255)"
	TypeTimestamp DataType = "timestamp"
	TypeDouble    DataType = "double"
)

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateTableReq struct {
	TableName string   `json:"name" binding:"required"`
	Columns   []Column `json:"columns" binding:"required"`
	SessionID string   `json:"sessionid"`
}

type InsertRowReq struct {
	TableName string         `json:"tablename"`
	Row       map[string]any `json:"row"`
}

type EditRowReq struct {
	TableName string         `json:"tablename"`
	RowId     int            `json:"rowid"`
	Row       map[string]any `json:"row"`
}

type DeleteRowReq struct {
	TableName string `json:"tablename"`
	RowId     int    `json:"rowid"`
}

type CreateSessionReq struct {
	TableName string `json:"tablename"`
}

func (req *CreateTableReq) ValidateNoDuplicateColumns() error {
	seen := make(map[string]bool)
	for _, col := range req.Columns {
		normalizedName := strings.ToLower(col.Name) // case-insensitive comparison
		if seen[normalizedName] {
			return fmt.Errorf("duplicate column name found: %s", col.Name)
		}
		seen[normalizedName] = true
	}
	return nil
}
