package machina

import (
	"context"
	"testing"
)

// MockCondition is a test condition implementation
func MockCondition(ctx context.Context, data map[string]any) (bool, error) {
	return true, nil
}

// MockAction is a test action implementation
func MockAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	return nil, nil
}

func TestRegistry_RegisterAndGetCondition(t *testing.T) {
	registry := NewRegistry()

	// Register condition
	err := registry.RegisterCondition("testCondition", MockCondition)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Get condition
	retrieved, err := registry.GetCondition("testCondition")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved == nil {
		t.Error("Expected condition function, got nil")
	}
}

func TestRegistry_RegisterConditionTwice(t *testing.T) {
	registry := NewRegistry()

	// Register condition
	err := registry.RegisterCondition("testCondition", MockCondition)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Try to register the same condition again
	err = registry.RegisterCondition("testCondition", MockCondition)
	if err == nil {
		t.Error("Expected error when registering condition twice, got nil")
	}
}

func TestRegistry_GetNonExistentCondition(t *testing.T) {
	registry := NewRegistry()

	// Try to get a non-existent condition
	_, err := registry.GetCondition("nonExistent")
	if err == nil {
		t.Error("Expected error when getting non-existent condition, got nil")
	}
}

func TestRegistry_RegisterAndGetAction(t *testing.T) {
	registry := NewRegistry()

	// Register action
	err := registry.RegisterAction("testAction", MockAction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Get action
	retrieved, err := registry.GetAction("testAction")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved == nil {
		t.Error("Expected action function, got nil")
	}
}

func TestRegistry_RegisterActionTwice(t *testing.T) {
	registry := NewRegistry()

	// Register action
	err := registry.RegisterAction("testAction", MockAction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Try to register the same action again
	err = registry.RegisterAction("testAction", MockAction)
	if err == nil {
		t.Error("Expected error when registering action twice, got nil")
	}
}

func TestRegistry_GetNonExistentAction(t *testing.T) {
	registry := NewRegistry()

	// Try to get a non-existent action
	_, err := registry.GetAction("nonExistent")
	if err == nil {
		t.Error("Expected error when getting non-existent action, got nil")
	}
}
