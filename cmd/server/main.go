package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/rahulpahuja/go-machina/internal/core"
	"github.com/rahulpahuja/go-machina/machina"
)

func main() {
	// Load workflow definition
	definition, err := machina.LoadWorkflowDefinition("configs/workflow.yaml")
	if err != nil {
		slog.Error("Failed to load workflow definition", "error", err)
		os.Exit(1)
	}

	// Set initial state if not already set
	if definition.InitialState == "" {
		definition.InitialState = "start"
	}

	// Create registry and register conditions/actions
	registry := machina.NewRegistry()

	// Register conditions
	registry.RegisterCondition("isUserValid", core.IsUserValidCondition)
	registry.RegisterCondition("isPaymentSuccess", core.IsPaymentSuccessCondition)

	// Register actions
	registry.RegisterAction("logStart", core.LogStartAction)
	registry.RegisterAction("logProcessing", core.LogProcessingAction)
	registry.RegisterAction("logCompletion", core.LogCompletionAction)
	registry.RegisterAction("logFailure", core.LogFailureAction)
	registry.RegisterAction("chargePayment", core.ChargePaymentAction)
	registry.RegisterAction("sendReceipt", core.SendReceiptAction)
	registry.RegisterAction("handleFailure", core.HandleFailureAction)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := machina.NewStateMachine(definition, registry, logger)

	// Execute workflow step by step
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderData := map[string]any{
		"orderId": "12345",
		"user": map[string]any{
			"id":    "user123",
			"email": "user@example.com",
		},
		"amount": 99.99,
	}

	slog.Info("Starting workflow execution")

	// Start with initial state from definition
	currentState := definition.InitialState

	// Execute workflow with auto-event handling
	for {
		// In a real application, you would determine the event based on external triggers
		// For this example, we'll use predefined events
		event := ""
		switch currentState {
		case "start":
			event = "validate"
		case "processOrder":
			// Simulate random success/failure for the process event
			event = "process" // In a real app, this might be determined by business logic
		default:
			// If we don't know what event to trigger, break the loop
			slog.Info("No event to trigger for current state", "state", currentState)
			break
		}

		if event == "" {
			break
		}

		// Trigger the event
		result, err := fsm.Trigger(ctx, currentState, event, orderData)
		if err != nil {
			slog.Error("Workflow execution failed", "error", err)
			os.Exit(1)
		}

		// Update state and data
		currentState = result.NewState
		orderData = result.PersistenceData

		slog.Info("Transition completed", "newState", result.NewState, "autoEvent", result.AutoEvent)

		// Handle auto events
		if result.AutoEvent != "" {
			slog.Info("Auto event triggered", "event", result.AutoEvent)
			continue // Continue the loop with the auto event
		}

		// Break if we've reached a terminal state or have no more events to process
		if currentState == "complete" || currentState == "failed" {
			break
		}
	}

	// Final state
	if currentState == "complete" {
		slog.Info("Workflow completed successfully", "result", orderData)
	} else if currentState == "failed" {
		slog.Info("Workflow failed", "result", orderData)
	} else {
		slog.Info("Workflow in progress", "current_state", currentState)
	}
}
