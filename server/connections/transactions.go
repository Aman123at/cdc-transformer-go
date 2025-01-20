package connections

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Aman123at/cdc-go/models"
	"github.com/Aman123at/cdc-go/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateNewTable(data models.CreateTableReq) error {
	conn, poolerr := PgPool.Get()
	defer PgPool.Put(conn)
	if poolerr != nil {
		log.Println("Unable to get connection from pool")
		return poolerr
	}

	query, err := utils.GenerateCreateTableQuery(data)
	if err != nil {
		return err
	}

	res, execerr := conn.dbInstance.Exec(query)
	if execerr != nil {
		return execerr
	}
	log.Println(res.LastInsertId())
	log.Println(res.RowsAffected())

	return nil
}

func InsertRow(rowData models.InsertRowReq) error {
	conn, poolerr := PgPool.Get()
	defer PgPool.Put(conn)
	if poolerr != nil {
		log.Println("Unable to get connection from pool")
		return poolerr
	}

	var columns []string
	var values []interface{}
	var placeholders []string
	i := 1

	for col, val := range rowData.Row {
		columns = append(columns, col)
		values = append(values, val)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}

	// Construct the INSERT query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		rowData.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute the query with values
	_, err := conn.dbInstance.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil

}

func EditRow(rowData models.EditRowReq, whereClause map[string]any) error {
	conn, poolerr := PgPool.Get()
	defer PgPool.Put(conn)
	if poolerr != nil {
		log.Println("Unable to get connection from pool")
		return poolerr
	}

	// Build SET clause
	var setValues []string
	var values []interface{}
	i := 1

	for col, val := range rowData.Row {
		setValues = append(setValues, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	// Build WHERE clause
	var whereConditions []string
	for col, val := range whereClause {
		whereConditions = append(whereConditions, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	// Construct the UPDATE query
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		rowData.TableName,
		strings.Join(setValues, ", "),
		strings.Join(whereConditions, " AND "),
	)

	// Execute the query with values
	_, err := conn.dbInstance.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRow(tablename string, whereClause map[string]any) error {
	conn, poolerr := PgPool.Get()
	defer PgPool.Put(conn)
	if poolerr != nil {
		log.Println("Unable to get connection from pool")
		return poolerr
	}

	// Build WHERE clause
	var whereConditions []string
	var values []interface{}
	i := 1

	for col, val := range whereClause {
		whereConditions = append(whereConditions, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	// Construct the DELETE query
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		tablename,
		strings.Join(whereConditions, " AND "),
	)

	// Execute the query with values
	_, err := conn.dbInstance.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

type Column struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Constraint interface{} `json:"constraint"`
}

type TableData struct {
	TableName string                   `json:"tablename"`
	Rows      []map[string]interface{} `json:"rows"`
	Columns   []Column                 `json:"columns"`
}

func GetAllTablesData() ([]TableData, error) {
	conn, poolerr := PgPool.Get()
	defer PgPool.Put(conn)
	if poolerr != nil {
		return nil, fmt.Errorf("error getting connection from pool: %v", poolerr)
	}

	rows, err := conn.dbInstance.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return nil, fmt.Errorf("error fetching tables: %v", err)
	}
	defer rows.Close()

	var allData []TableData

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("error scanning table name: %v", err)
		}
		if tableName != "ddl_changes" {

			columnQuery := `
                SELECT 
                    c.column_name,
                    c.data_type,
                    CASE 
                        WHEN c.is_nullable = 'NO' THEN 'NOT NULL'
                        ELSE 'NULL'
                    END as nullable,
                    CASE 
                        WHEN c.column_default IS NOT NULL THEN c.column_default
                        ELSE ''
                    END as default_value,
                    CASE 
                        WHEN pk.column_name IS NOT NULL THEN true
                        ELSE false
                    END as is_primary_key
                FROM information_schema.columns c
                LEFT JOIN (
                    SELECT ku.column_name
                    FROM information_schema.table_constraints tc
                    JOIN information_schema.key_column_usage ku
                        ON tc.constraint_name = ku.constraint_name
                    WHERE tc.constraint_type = 'PRIMARY KEY'
                        AND tc.table_name = $1
                ) pk ON c.column_name = pk.column_name
                WHERE c.table_name = $1
                ORDER BY c.ordinal_position;`

			columnRows, err := conn.dbInstance.Query(columnQuery, tableName)
			if err != nil {
				return nil, fmt.Errorf("error fetching column info for table %s: %v", tableName, err)
			}
			defer columnRows.Close()

			var columns []Column
			for columnRows.Next() {
				var (
					colName, dataType, nullable, defaultValue string
					isPrimaryKey                              bool
				)

				if err := columnRows.Scan(&colName, &dataType, &nullable, &defaultValue, &isPrimaryKey); err != nil {
					return nil, fmt.Errorf("error scanning column info: %v", err)
				}

				constraints := make(map[string]interface{})
				constraints["nullable"] = nullable
				if defaultValue != "" {
					constraints["default"] = defaultValue
				}
				if isPrimaryKey {
					constraints["primary_key"] = true
				}

				columns = append(columns, Column{
					Name:       colName,
					Type:       dataType,
					Constraint: constraints,
				})
			}

			// Get table rows
			tableRows, err := conn.dbInstance.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
			if err != nil {
				return nil, fmt.Errorf("error fetching rows from table %s: %v", tableName, err)
			}
			defer tableRows.Close()

			tableData := TableData{
				TableName: tableName,
				Columns:   columns,
				Rows:      make([]map[string]interface{}, 0),
			}

			// Get column names for row data
			colNames, err := tableRows.Columns()
			if err != nil {
				return nil, fmt.Errorf("error getting columns for table %s: %v", tableName, err)
			}

			// ... existing row scanning code ...
			for tableRows.Next() {
				values := make([]interface{}, len(colNames))
				valuePtrs := make([]interface{}, len(colNames))
				for i := range colNames {
					valuePtrs[i] = &values[i]
				}

				if err := tableRows.Scan(valuePtrs...); err != nil {
					return nil, fmt.Errorf("error scanning row: %v", err)
				}

				row := make(map[string]interface{})
				for i, col := range colNames {
					var v interface{}
					val := values[i]
					b, ok := val.([]byte)
					if ok {
						v = string(b)
					} else {
						v = val
					}
					row[col] = v
				}

				tableData.Rows = append(tableData.Rows, row)
			}

			allData = append(allData, tableData)
		}
	}

	return allData, nil
}

type MongoCollectionData struct {
	CollectionName string                   `json:"collectionname"`
	Documents      []map[string]interface{} `json:"documents"`
}

func GetAllCollectionsData() ([]MongoCollectionData, error) {
	var allData []MongoCollectionData

	// Get all collection names
	collections, err := MongodbConn.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching collections: %v", err)
	}

	// Iterate through each collection
	for _, collName := range collections {
		collection := MongodbConn.Collection(collName)

		// Find all documents in the collection
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, fmt.Errorf("error fetching documents from collection %s: %v", collName, err)
		}
		defer cursor.Close(context.Background())

		// Create collection data structure
		collData := MongoCollectionData{
			CollectionName: collName,
			Documents:      make([]map[string]interface{}, 0),
		}

		// Iterate through all documents
		for cursor.Next(context.Background()) {
			var document map[string]interface{}
			if err := cursor.Decode(&document); err != nil {
				return nil, fmt.Errorf("error decoding document: %v", err)
			}
			collData.Documents = append(collData.Documents, document)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %v", err)
		}

		allData = append(allData, collData)
	}

	return allData, nil
}
