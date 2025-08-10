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

// RecordPreviousStateAction records the previous state before entering a side quest
// and pushes it onto the workflow stack
func RecordPreviousStateAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	previousState := data["state"]
	fmt.Printf("Recording previous state before side quest: %s\n", previousState)
	
	// Get the current workflow stack or create a new one
	var workflowStack []string
	if stack, ok := data["WorkflowStack"].([]string); ok {
		workflowStack = stack
	}
	
	// Push the previous state onto the stack
	workflowStack = append(workflowStack, previousState.(string))
	
	return map[string]any{
		"WorkflowStack": workflowStack,
	}, nil
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
	registry.RegisterAction("recordPreviousState", RecordPreviousStateAction)

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
	data := map[string]any{"state": currentState}

	fmt.Println("Starting dynamic workflow with side quests")
	fmt.Println("Main flow: A -> B -> C -> D -> E -> F -> G")
	fmt.Println("Side quests: B# and C#")
	
	// Normal flow to C
	result, err := fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	result, err = fsm.Trigger(ctx, currentState, "next", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Now take a side quest from C to B#
	fmt.Println("\n--- Taking side quest from C to B# ---")
	result, err = fsm.Trigger(ctx, currentState, "sideQuestB", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Return from B# to previous state (C) using dynamic transition
	fmt.Println("\n--- Returning from B# to previous state ---")
	result, err = fsm.Trigger(ctx, currentState, "return", data)
	if err != nil {
		fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
		return
	}
	currentState = result.NewState
	data = result.PersistenceData
	data["state"] = currentState

	// Continue normal flow to G
	fmt.Println("\n--- Continuing normal flow to G ---")
	for currentState != "G" {
		result, err := fsm.Trigger(ctx, currentState, "next", data)
		if err != nil {
			fmt.Printf("Error transitioning from %s: %v\n", currentState, err)
			return
		}
		currentState = result.NewState
		data = result.PersistenceData
		data["state"] = currentState
	}

	fmt.Printf("Workflow completed. Final state: %s\n", currentState)
	fmt.Println("\nNote: This implementation uses the dynamic transition features")
	fmt.Println("of GoMachina, with the __RETURN_TO_PREVIOUS_STATE__ action handling")
	fmt.Println("the return from side quests automatically.")
}