package connections

import (
	"context"
	"testing"
)

func TestCheckDBConnections(t *testing.T) {
	// Save original connection values
	origPgConn := PgConn
	origMongoConn := MongodbConn

	defer func() {
		// Restore original connections after test
		PgConn = origPgConn
		MongodbConn = origMongoConn
	}()

	// Test with invalid connection strings (should panic)
	t.Run("test connections", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with invalid connection strings")
			}
		}()

		CheckDBConnetions()
	})
}

func TestMongoConnection(t *testing.T) {
	if MongodbConn != nil {
		err := MongodbConn.Client().Ping(context.Background(), nil)
		if err != nil {
			t.Errorf("MongoDB connection is not alive: %v", err)
		}
	}
}

func TestPgConnection(t *testing.T) {
	if PgConn != nil {
		err := PgConn.Ping(context.Background())
		if err != nil {
			t.Errorf("PostgreSQL connection is not alive: %v", err)
		}
	}
}
