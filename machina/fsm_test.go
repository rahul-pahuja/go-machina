package machina

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestStateMachine_Trigger(t *testing.T) {
	// Create a simple workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysTrue",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysTrue", MockTrueCondition)

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Trigger event
	ctx := context.Background()
	result, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.NewState != "end" {
		t.Errorf("Expected new state to be 'end', got %s", result.NewState)
	}

	if result.AutoEvent != "" {
		t.Errorf("Expected no auto event, got %s", result.AutoEvent)
	}

	if result.PersistenceData == nil {
		t.Error("Expected persistence data, got nil")
	}
}

func TestStateMachine_Trigger_ConditionFailure(t *testing.T) {
	// Create a workflow definition with a failing condition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysFalse",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysFalse", MockFalseCondition)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event
	ctx := context.Background()
	_, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStateMachine_Trigger_ActionExecution(t *testing.T) {
	// Create a workflow definition with actions
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysTrue",
						},
						Actions: []string{
							"updateAction",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysTrue", MockTrueCondition)
	registry.RegisterAction("updateAction", MockUpdateAction)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event
	ctx := context.Background()
	result, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.NewState != "end" {
		t.Errorf("Expected new state to be 'end', got %s", result.NewState)
	}

	if result.AutoEvent != "" {
		t.Errorf("Expected no auto event, got %s", result.AutoEvent)
	}

	// Check if the action updated the data
	if updated, exists := result.PersistenceData["updated"]; !exists || !updated.(bool) {
		t.Error("Expected persistence data to be updated by action")
	}
}

func TestStateMachine_Trigger_ConditionError(t *testing.T) {
	// Create a workflow definition with a failing condition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"errorCondition",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("errorCondition", MockErrorCondition)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event
	ctx := context.Background()
	_, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStateMachine_Trigger_ActionError(t *testing.T) {
	// Create a workflow definition with a failing action
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysTrue",
						},
						Actions: []string{
							"errorAction",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysTrue", MockTrueCondition)
	registry.RegisterAction("errorAction", MockErrorAction)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event
	ctx := context.Background()
	_, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStateMachine_Trigger_ContextTimeout(t *testing.T) {
	// Create a workflow definition with a slow action
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysTrue",
						},
						Actions: []string{
							"slowAction",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysTrue", MockTrueCondition)
	registry.RegisterAction("slowAction", MockSlowAction)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestStateMachine_Trigger_LifecycleHooks(t *testing.T) {
	// Create a workflow definition with lifecycle hooks
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				OnLeave: []string{
					"onLeaveAction",
				},
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"alwaysTrue",
						},
						Actions: []string{
							"transitionAction",
						},
					},
				},
			},
			"end": {
				Name: "end",
				OnEnter: []string{
					"onEnterAction",
				},
			},
		},
	}

	// Create registry with mock implementations
	registry := NewRegistry()
	registry.RegisterCondition("alwaysTrue", MockTrueCondition)
	registry.RegisterAction("onLeaveAction", MockUpdateAction)
	registry.RegisterAction("transitionAction", MockUpdateAction)
	registry.RegisterAction("onEnterAction", MockUpdateAction)

	// Create state machine
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fsm := NewStateMachine(definition, registry, logger)

	// Trigger event
	ctx := context.Background()
	result, err := fsm.Trigger(ctx, "start", "proceed", map[string]any{})

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.NewState != "end" {
		t.Errorf("Expected new state to be 'end', got %s", result.NewState)
	}

	if result.AutoEvent != "" {
		t.Errorf("Expected no auto event, got %s", result.AutoEvent)
	}

	if result.PersistenceData == nil {
		t.Error("Expected persistence data, got nil")
	}
}
