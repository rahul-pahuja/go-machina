package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/rahulpahuja/go-machina/machina"
)

// LogAction logs the current state
func LogAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	state := data["state"]
	fmt.Printf("Entering state: %s\n", state)
	return nil, nil
}

// IsEvenCondition checks if a number in the data is even
func IsEvenCondition(ctx context.Context, data map[string]any) (bool, error) {
	if num, ok := data["number"]; ok {
		if n, ok := num.(int); ok {
			return n%2 == 0, nil
		}
	}
	return false, nil
}

// IsOddCondition checks if a number in the data is odd
func IsOddCondition(ctx context.Context, data map[string]any) (bool, error) {
	if num, ok := data["number"]; ok {
		if n, ok := num.(int); ok {
			return n%2 != 0, nil
		}
	}
	return false, nil
}

func main() {
	// Load workflow definition from YAML file
	definition, err := machina.LoadWorkflowDefinition("workflow.yaml")
	if err != nil {
		fmt.Printf("Error loading workflow definition: %v\n", err)
		return
	}

	// Create registry and register actions and conditions
	registry := machina.NewRegistry()
	registry.RegisterAction("logAction", LogAction)
	registry.RegisterCondition("isEven", IsEvenCondition)
	registry.RegisterCondition("isOdd", IsOddCondition)

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create state machine
	fsm := machina.NewStateMachine(definition, registry, logger)
	if fsm == nil {
		fmt.Println("Failed to create state machine")
		return
	}

	// Execute the workflow with an even number (Branch 1: A -> B -> D -> E)
	ctx := context.Background()
	currentState := "A"
	data := map[string]any{"number": 4, "state": currentState}

	fmt.Println("Starting conditional workflow with even number (4): A -> B -> D -> E")
	
	// Transition from A to B (using process event with isEven condition)
	result, err := fsm.Trigger(ctx, currentState, "process", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Transition from B to D
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Transition from D to E
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	fmt.Printf("Workflow completed with even number. Final state: %s\n", currentState)

	// Execute the workflow with an odd number (Branch 2: A -> C -> E)
	fmt.Println("\n==================================================")
	currentState = "A"
	data = map[string]any{"number": 7, "state": currentState}

	fmt.Println("Starting conditional workflow with odd number (7): A -> C -> E")
	
	// Transition from A to C (using process event with isOdd condition)
	result, err = fsm.Trigger(ctx, currentState, "process", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Transition from C to E
	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	fmt.Printf("Workflow completed with odd number. Final state: %s\n", currentState)
}