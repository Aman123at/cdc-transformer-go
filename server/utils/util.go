package utils

import (
	"fmt"
	"strings"

	"github.com/Aman123at/cdc-go/models"
)

// GenerateCreateTableQuery generates a CREATE TABLE SQL query from the given request
func GenerateCreateTableQuery(data models.CreateTableReq) (string, error) {

	var columnDefs []string
	for _, col := range data.Columns {
		name := col.Name
		typeStr := col.Type

		// Validate that the type is one of our defined DataTypes
		isValid := false
		for _, validType := range []models.DataType{models.TypeInt, models.TypeText, models.TypeBoolean, models.TypeVarchar, models.TypeTimestamp, models.TypeDouble} {
			if models.DataType(typeStr) == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return "", fmt.Errorf("invalid data type: %s", typeStr)
		}

		if name == "id" {
			typeStr = "serial primary key"
		}

		// Build the column definition
		columnDef := fmt.Sprintf("%s %s", strings.ToLower(name), typeStr)
		columnDefs = append(columnDefs, columnDef)
	}

	tablename := strings.ToLower(data.TableName)
	// Construct the final query
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n    %s\n);",
		tablename,
		strings.Join(columnDefs, ",\n    "),
	)

	return query, nil
}

func TrimColumnSpaces(data models.Column) {
	data.Name = strings.TrimSpace(data.Name)
}
