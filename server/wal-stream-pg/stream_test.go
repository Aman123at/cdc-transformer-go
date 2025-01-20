package stream

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Aman123at/cdc-go/connections"
	"github.com/jackc/pglogrepl"
)

func TestStartStream(t *testing.T) {
	// This is more of an integration test
	// We'll test it with a small timeout to avoid hanging
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StartStream panicked: %v", r)
			}
			done <- true
		}()

		// Only run for a short duration
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go StartStream()

		<-ctx.Done()
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(3 * time.Second):
		t.Error("Test timed out")
	}
}

func TestReplicationSlotCreation(t *testing.T) {
	if connections.PgConn == nil {
		t.Skip("PostgreSQL connection not available")
	}

	slotName := "test_slot"
	query := fmt.Sprintf("SELECT pg_drop_replication_slot(%s)", slotName)
	// Clean up any existing slot
	connections.PgConn.Exec(context.Background(), query)

	// Test creating a new slot
	_, err := pglogrepl.CreateReplicationSlot(
		context.Background(),
		connections.PgConn,
		slotName,
		"wal2json",
		pglogrepl.CreateReplicationSlotOptions{Temporary: true},
	)

	if err != nil {
		t.Errorf("Failed to create replication slot: %v", err)
	}
}
