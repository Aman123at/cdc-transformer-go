package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WALEvent represents the structure of a WAL change event
type WALEvent struct {
	Change []ChangeEvent `json:"change"`
}

type ChangeEvent struct {
	Kind         string   `json:"kind"`
	Schema       string   `json:"schema"`
	Table        string   `json:"table"`
	ColumnNames  []string `json:"columnnames"`
	ColumnTypes  []string `json:"columntypes"`
	ColumnValues []any    `json:"columnvalues"`
	OldKeys      *KeyInfo `json:"oldkeys,omitempty"`
}

type KeyInfo struct {
	KeyNames  []string `json:"keynames"`
	KeyTypes  []string `json:"keytypes"`
	KeyValues []any    `json:"keyvalues"`
}

type WALParser struct {
	db *mongo.Database
}

func NewWALParser(db *mongo.Database) *WALParser {
	return &WALParser{db: db}
}

func (wp *WALParser) ProcessWALEvent(event WALEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if wp.db == nil {
		log.Println("MongoDB Connection is not available")
		return nil
	}

	for _, change := range event.Change {
		if change.Table == "ddl_changes" {
			continue
		}
		collection := wp.db.Collection(change.Table)

		switch change.Kind {
		case "insert":
			if err := wp.handleInsert(ctx, collection, change); err != nil {
				return err
			}

		case "update":
			if err := wp.handleUpdate(ctx, collection, change); err != nil {
				return err
			}

		case "delete":
			if err := wp.handleDelete(ctx, collection, change); err != nil {
				return err
			}
		}
	}

	return nil
}

func (wp *WALParser) handleInsert(ctx context.Context, collection *mongo.Collection, change ChangeEvent) error {
	doc := wp.createDocument(change.ColumnNames, change.ColumnValues)
	_, err := collection.InsertOne(ctx, doc)
	return err
}

func (wp *WALParser) handleUpdate(ctx context.Context, collection *mongo.Collection, change ChangeEvent) error {
	if change.OldKeys == nil {
		return fmt.Errorf("no old keys provided for update operation")
	}

	filter := wp.createFilterFromKeys(change.OldKeys)
	update := bson.M{"$set": wp.createDocument(change.ColumnNames, change.ColumnValues)}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (wp *WALParser) handleDelete(ctx context.Context, collection *mongo.Collection, change ChangeEvent) error {
	if change.OldKeys == nil {
		return fmt.Errorf("no old keys provided for delete operation")
	}

	filter := wp.createFilterFromKeys(change.OldKeys)
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

func (wp *WALParser) createDocument(columnNames []string, columnValues []any) bson.M {
	doc := bson.M{}
	for i := range columnNames {
		doc[columnNames[i]] = columnValues[i]
	}
	return doc
}

func (wp *WALParser) createFilterFromKeys(keys *KeyInfo) bson.M {
	filter := bson.M{}
	for i := range keys.KeyNames {
		filter[keys.KeyNames[i]] = keys.KeyValues[i]
	}
	return filter
}
