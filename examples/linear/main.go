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

func main() {
	// Load workflow definition from YAML file
	definition, err := machina.LoadWorkflowDefinition("workflow.yaml")
	if err != nil {
		fmt.Printf("Error loading workflow definition: %v\n", err)
		return
	}

	// Create registry and register actions
	registry := machina.NewRegistry()
	registry.RegisterAction("logAction", LogAction)

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create state machine
	fsm := machina.NewStateMachine(definition, registry, logger)
	if fsm == nil {
		fmt.Println("Failed to create state machine")
		return
	}

	// Execute the workflow
	ctx := context.Background()
	currentState := "A"

	fmt.Println("Starting simple linear workflow: A -> B -> C")
	
	// Transition from A to B
	result, err := fsm.Trigger(ctx, currentState, "next", map[string]any{"state": currentState})
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	// Transition from B to C
	result, err = fsm.Trigger(ctx, currentState, "next", map[string]any{"state": currentState})
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState

	fmt.Printf("Workflow completed. Final state: %s\n", currentState)
}