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

	// Start state
	currentState := "start"

	// Trigger validate event
	newState, result, err := fsm.Trigger(ctx, currentState, "validate", orderData)
	if err != nil {
		slog.Error("Workflow execution failed", "error", err)
		os.Exit(1)
	}
	currentState = newState
	orderData = result

	// Trigger process event (with 80% success rate)
	newState, result, err = fsm.Trigger(ctx, currentState, "process", orderData)
	if err != nil {
		slog.Error("Workflow execution failed", "error", err)
		os.Exit(1)
	}
	currentState = newState
	orderData = result

	// If we reached the complete state, the workflow is done
	if currentState == "complete" {
		slog.Info("Workflow completed successfully", "result", orderData)
	} else if currentState == "failed" {
		slog.Info("Workflow failed", "result", orderData)
	} else {
		slog.Info("Workflow in progress", "current_state", currentState)
	}
}
