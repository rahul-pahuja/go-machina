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

func main() {
	// Load workflow definition from YAML file
	definition, err := machina.LoadWorkflowDefinition("workflow_skip.yaml")
	if err != nil {
		fmt.Printf("Error loading workflow definition: %v\n", err)
		return
	}

	// Create registry and register actions
	registry := machina.NewRegistry()
	registry.RegisterAction("logAction", LogAction)
	registry.RegisterAction("timerAction", TimerAction)

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

	fmt.Println("Starting timeout workflow demonstration")
	fmt.Println("Normal flow: A -> B -> C -> D -> E")
	
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

	// Simulate timeout scenario - if timer expires and we're in state C,
	// we should reset but skip B in the retry
	fmt.Println("\n--- Simulating timeout scenario ---")
	fmt.Println("Timer expired while in state C. Resetting workflow.")
	fmt.Println("Retry path should be: A -> C -> D (skipping B)")

	// Reset to A
	currentState = "A"
	data = map[string]any{"state": currentState, "retry": true}

	// Transition from A to B again (but this time it's a retry)
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// In a real implementation, we would check if this is a retry
	// and skip B, going directly to C. But for this example, we'll
	// simulate that by manually setting the state to C.
	fmt.Println("Skipping state B in retry (simulated)")
	currentState = "C"
	data["state"] = currentState

	// Continue from C
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Continue to E
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	fmt.Printf("Workflow completed. Final state: %s\n", currentState)
	fmt.Println("\nNote: In a real implementation, the library would need")
	fmt.Println("additional logic to detect retry scenarios and skip states.")
}