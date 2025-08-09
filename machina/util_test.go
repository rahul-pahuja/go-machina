package machina

import (
	"testing"
)

func TestMergeData(t *testing.T) {
	// Create a state machine
	sm := &StateMachine{}

	// Test merging two maps
	original := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	updates := map[string]any{
		"key2": "updated_value2",
		"key3": "value3",
	}

	result := sm.mergeData(original, updates)

	// Verify the result
	if result["key1"] != "value1" {
		t.Errorf("Expected key1 to be 'value1', got '%v'", result["key1"])
	}

	if result["key2"] != "updated_value2" {
		t.Errorf("Expected key2 to be 'updated_value2', got '%v'", result["key2"])
	}

	if result["key3"] != "value3" {
		t.Errorf("Expected key3 to be 'value3', got '%v'", result["key3"])
	}
}

func TestMergeData_EmptyOriginal(t *testing.T) {
	// Create a state machine
	sm := &StateMachine{}

	// Test merging with empty original
	original := map[string]any{}

	updates := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	result := sm.mergeData(original, updates)

	// Verify the result
	if result["key1"] != "value1" {
		t.Errorf("Expected key1 to be 'value1', got '%v'", result["key1"])
	}

	if result["key2"] != "value2" {
		t.Errorf("Expected key2 to be 'value2', got '%v'", result["key2"])
	}
}

func TestMergeData_EmptyUpdates(t *testing.T) {
	// Create a state machine
	sm := &StateMachine{}

	// Test merging with empty updates
	original := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	updates := map[string]any{}

	result := sm.mergeData(original, updates)

	// Verify the result
	if result["key1"] != "value1" {
		t.Errorf("Expected key1 to be 'value1', got '%v'", result["key1"])
	}

	if result["key2"] != "value2" {
		t.Errorf("Expected key2 to be 'value2', got '%v'", result["key2"])
	}
}

func TestMergeData_BothEmpty(t *testing.T) {
	// Create a state machine
	sm := &StateMachine{}

	// Test merging two empty maps
	original := map[string]any{}
	updates := map[string]any{}

	result := sm.mergeData(original, updates)

	// Verify the result is empty
	if len(result) != 0 {
		t.Errorf("Expected empty map, got map with %d elements", len(result))
	}
}

func TestRecordTransitionError(t *testing.T) {
	// Create a state machine without metrics
	sm := &StateMachine{}

	// This should not panic or cause any issues
	sm.recordTransitionError("from", "event", "error_type", nil)

	// Test with metrics
	// Create a test registry
	// reg := prometheus.NewRegistry()
	// metrics := NewMetrics(reg)
	// smWithMetrics := &StateMachine{
	// 	metrics: metrics,
	// }
	//
	// This should record the error
	// smWithMetrics.recordTransitionError("from", "event", "error_type", nil)
}
