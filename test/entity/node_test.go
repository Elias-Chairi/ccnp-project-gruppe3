package entity_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/entity"
)

// TestSetID tests the SetID method of the Node struct.
func TestSetID(t *testing.T) {
	node := entity.NewNode()
	err := node.SetID("node-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if node.ID != "node-1" {
		t.Fatalf("expected ID to be 'node-1', got %v", node.ID)
	}
}

// TestSetID_Empty tests the SetID method with an empty string, expecting an error.
func TestSetID_Empty(t *testing.T) {
	node := entity.NewNode()
	err := node.SetID("")
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
	if node.ID != "" {
		t.Fatalf("expected ID to be '', got %v", node.ID)
	}
}
