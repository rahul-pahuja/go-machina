package main

import (
	"context"
	"testing"
	"time"

	"github.com/rahulpahuja/go-machina/machina"
)

// ExpiryCondition simulates a condition that checks if a process has expired
func ExpiryCondition(ctx context.Context, data map[string]any) (bool, error) {
	expiryTime, ok := data["expiryTime"].(time.Time)
	if !ok {
		// Default to not expired
		return false, nil
	}
	return time.Now().After(expiryTime), nil
}

// NotExpiredCondition simulates a condition that checks if a process has not expired
func NotExpiredCondition(ctx context.Context, data map[string]any) (bool, error) {
	expiryTime, ok := data["expiryTime"].(time.Time)
	if !ok {
		// Default to not expired
		return true, nil
	}
	return time.Now().Before(expiryTime), nil
}

// TimeoutAction simulates an action that takes a long time to execute
func TimeoutAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	duration, ok := data["duration"].(time.Duration)
	if !ok {
		duration = 100 * time.Millisecond
	}

	select {
	case <-time.After(duration):
		return nil, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// SuccessAction simulates a successful action
func SuccessAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	return map[string]any{
		"success": true,
	}, nil
}

// UpdateAction simulates an action that updates data
func UpdateAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	result := make(map[string]any)

	// Copy existing data
	for k, v := range data {
		result[k] = v
	}

	// Add/update with new data
	result["updated"] = true
	result["timestamp"] = time.Now()

	return result, nil
}

func TestIntegration_ExpiryScenario(t *testing.T) {
	// Create a workflow definition with expiry scenario
	definition := &machina.WorkflowDefinition{
		States: map[string]machina.State{
			"start": {
				Name: "start",
				Transitions: []machina.Transition{
					{
						Event:  "check",
						Target: "process",
						Conditions: []string{
							"notExpired",
						},
					},
					{
						Event:  "check",
						Target: "timeout",
						Conditions: []string{
							"expired",
						},
					},
				},
			},
			"process": {
				Name: "process",
				Transitions: []machina.Transition{
					{
						Event:  "complete",
						Target: "success",
						Actions: []string{
							"updateAction",
						},
					},
				},
			},
			"timeout": {
				Name: "timeout",
				OnEnter: []string{
					"successAction",
				},
			},
			"success": {
				Name: "success",
				OnEnter: []string{
					"successAction",
				},
			},
		},
	}

	// Create registry and register conditions/actions
	registry := machina.NewRegistry()

	// Register conditions
	registry.RegisterCondition("notExpired", NotExpiredCondition)
	registry.RegisterCondition("expired", ExpiryCondition)

	// Register actions
	registry.RegisterAction("updateAction", UpdateAction)
	registry.RegisterAction("successAction", SuccessAction)

	// Create state machine
	fsm := machina.NewStateMachine(definition, registry, nil)

	// Execute workflow
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := map[string]any{
		"processId":  "12345",
		"expiryTime": time.Now().Add(1 * time.Hour), // Not expired
	}

	// Trigger the check event
	currentState := "start"
	newState, result, err := fsm.Trigger(ctx, currentState, "check", data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if newState != "process" {
		t.Errorf("Expected new state to be 'process', got %s", newState)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	// Trigger the complete event
	currentState = newState
	newState, result, err = fsm.Trigger(ctx, currentState, "complete", result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if newState != "success" {
		t.Errorf("Expected new state to be 'success', got %s", newState)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestIntegration_TimeoutScenario(t *testing.T) {
	// Create a workflow definition with timeout scenario
	definition := &machina.WorkflowDefinition{
		States: map[string]machina.State{
			"start": {
				Name: "start",
				Transitions: []machina.Transition{
					{
						Event:  "process",
						Target: "working",
						Actions: []string{
							"timeoutAction",
						},
					},
				},
			},
			"working": {
				Name: "working",
			},
		},
	}

	// Create registry and register conditions/actions
	registry := machina.NewRegistry()

	// Register actions
	registry.RegisterAction("timeoutAction", TimeoutAction)

	// Create state machine
	fsm := machina.NewStateMachine(definition, registry, nil)

	// Execute workflow with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	data := map[string]any{
		"processId": "12345",
		"duration":  1 * time.Second, // Action will take longer than context timeout
	}

	// Trigger the process event
	_, _, err := fsm.Trigger(ctx, "start", "process", data)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
