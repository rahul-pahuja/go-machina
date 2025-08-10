package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/rahulpahuja/go-machina/machina"
)

// LogAction logs the current state
func LogAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	state := data["state"]
	fmt.Printf("Entering state: %s\n", state)
	return nil, nil
}

// TimerAction starts a timer
func TimerAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Timer started after state A")
	return map[string]any{"timerStarted": true, "timerStart": time.Now()}, nil
}

// ResetAction resets the workflow to the beginning
func ResetAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Timeout occurred! Resetting workflow to beginning.")
	return map[string]any{"reset": true}, nil
}

func main() {
	// Load workflow definition from YAML file
	definition, err := machina.LoadWorkflowDefinition("workflow_reset.yaml")
	if err != nil {
		fmt.Printf("Error loading workflow definition: %v\n", err)
		return
	}

	// Create registry and register actions and conditions
	registry := machina.NewRegistry()
	registry.RegisterAction("logAction", LogAction)
	registry.RegisterAction("timerAction", TimerAction)
	registry.RegisterAction("resetAction", ResetAction)
	
	// Simple condition that always returns true
	alwaysTrue := func(ctx context.Context, data map[string]any) (bool, error) {
		return true, nil
	}
	registry.RegisterCondition("alwaysTrue", alwaysTrue)

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create state machine
	fsm := machina.NewStateMachine(definition, registry, logger)
	if fsm == nil {
		fmt.Println("Failed to create state machine")
		return
	}

	// Execute the workflow normally first
	ctx := context.Background()
	currentState := "A"
	data := map[string]any{"state": currentState}

	fmt.Println("Starting timeout workflow with auto-reset")
	fmt.Println("Normal flow: A -> B -> C -> D -> E")
	fmt.Println("On timeout: Any state -> A (reset)")
	
	// Transition from A to B (timer starts)
	result, err := fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Transition from B to C
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	fmt.Printf("Current state: %s (timer started, but not expired yet)\n", currentState)

	// Simulate timeout - transition to A with timeout event
	fmt.Println("\n--- Simulating timeout ---")
	result, err = fsm.Trigger(ctx, currentState, "timeout", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	fmt.Printf("After timeout, workflow reset to state: %s\n", currentState)

	// Now continue normally again to complete the workflow: A -> B -> C -> D -> E
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// B -> C
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// C -> D
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// D -> E (final state)
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	fmt.Printf("Workflow completed. Final state: %s\n", currentState)
}