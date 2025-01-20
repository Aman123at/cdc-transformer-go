package services

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewWALParser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("creates new parser", func(mt *mtest.T) {
		parser := NewWALParser(mt.DB)
		if parser == nil {
			t.Error("Expected non-nil parser")
		}
		if parser.db != mt.DB {
			t.Error("Expected parser to have correct database reference")
		}
	})
}

func TestProcessWALEvent(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("handles nil database", func(mt *mtest.T) {
		parser := &WALParser{db: nil}
		event := WALEvent{
			Change: []ChangeEvent{
				{
					Kind:  "insert",
					Table: "test_table",
				},
			},
		}
		err := parser.ProcessWALEvent(event)
		if err != nil {
			t.Errorf("Expected nil error for nil database, got %v", err)
		}
	})

	mt.Run("processes insert event", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		parser := NewWALParser(mt.DB)
		event := WALEvent{
			Change: []ChangeEvent{
				{
					Kind:         "insert",
					Table:        "test_table",
					ColumnNames:  []string{"id", "name"},
					ColumnValues: []any{1, "test"},
				},
			},
		}

		err := parser.ProcessWALEvent(event)
		if err != nil {
			t.Errorf("Failed to process insert event: %v", err)
		}
	})
}

func TestHandleInsert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("successful insert", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		parser := NewWALParser(mt.DB)
		change := ChangeEvent{
			ColumnNames:  []string{"id", "name"},
			ColumnValues: []any{1, "test"},
		}

		err := parser.handleInsert(context.Background(), mt.Coll, change)
		if err != nil {
			t.Errorf("Failed to handle insert: %v", err)
		}
	})
}

func TestHandleUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("update without old keys", func(mt *mtest.T) {
		parser := NewWALParser(mt.DB)
		change := ChangeEvent{
			ColumnNames:  []string{"id", "name"},
			ColumnValues: []any{1, "test"},
		}

		err := parser.handleUpdate(context.Background(), mt.Coll, change)
		if err == nil {
			t.Error("Expected error for update without old keys")
		}
	})

	mt.Run("successful update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		parser := NewWALParser(mt.DB)
		change := ChangeEvent{
			ColumnNames:  []string{"id", "name"},
			ColumnValues: []any{1, "test"},
			OldKeys: &KeyInfo{
				KeyNames:  []string{"id"},
				KeyValues: []any{1},
			},
		}

		err := parser.handleUpdate(context.Background(), mt.Coll, change)
		if err != nil {
			t.Errorf("Failed to handle update: %v", err)
		}
	})
}

func TestHandleDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("delete without old keys", func(mt *mtest.T) {
		parser := NewWALParser(mt.DB)
		change := ChangeEvent{}

		err := parser.handleDelete(context.Background(), mt.Coll, change)
		if err == nil {
			t.Error("Expected error for delete without old keys")
		}
	})

	mt.Run("successful delete", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		parser := NewWALParser(mt.DB)
		change := ChangeEvent{
			OldKeys: &KeyInfo{
				KeyNames:  []string{"id"},
				KeyValues: []any{1},
			},
		}

		err := parser.handleDelete(context.Background(), mt.Coll, change)
		if err != nil {
			t.Errorf("Failed to handle delete: %v", err)
		}
	})
}

func TestCreateDocument(t *testing.T) {
	parser := &WALParser{}

	tests := []struct {
		name         string
		columnNames  []string
		columnValues []any
		want         int
	}{
		{
			name:         "empty document",
			columnNames:  []string{},
			columnValues: []any{},
			want:         0,
		},
		{
			name:         "single field",
			columnNames:  []string{"id"},
			columnValues: []any{1},
			want:         1,
		},
		{
			name:         "multiple fields",
			columnNames:  []string{"id", "name", "age"},
			columnValues: []any{1, "test", 25},
			want:         3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parser.createDocument(tt.columnNames, tt.columnValues)
			if len(doc) != tt.want {
				t.Errorf("createDocument() got %v fields, want %v", len(doc), tt.want)
			}
		})
	}
}
