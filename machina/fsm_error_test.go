package machina

import (
	"context"
	"testing"
)

func TestStateMachine_Trigger_StateNotFound(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event from non-existent state
	_, _, err := fsm.Trigger(context.Background(), "nonexistent", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent state, got nil")
	}
}

func TestStateMachine_Trigger_TransitionNotFound(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger non-existent event
	_, _, err := fsm.Trigger(context.Background(), "start", "nonexistent", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent transition, got nil")
	}
}

func TestStateMachine_Trigger_ConditionNotFound(t *testing.T) {
	// Create a workflow definition with a condition that doesn't exist
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"nonexistentCondition",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, _, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent condition, got nil")
	}
}

func TestStateMachine_Trigger_ActionNotFound(t *testing.T) {
	// Create a workflow definition with an action that doesn't exist
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Actions: []string{
							"nonexistentAction",
						},
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, _, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent action, got nil")
	}
}

func TestStateMachine_Trigger_TargetStateNotFound(t *testing.T) {
	// Create a workflow definition with a transition to a non-existent state
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "nonexistent",
					},
				},
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, _, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent target state, got nil")
	}
}

func TestStateMachine_Trigger_OnLeaveActionNotFound(t *testing.T) {
	// Create a workflow definition with an OnLeave action that doesn't exist
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				OnLeave: []string{
					"nonexistentAction",
				},
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, _, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent OnLeave action, got nil")
	}
}

func TestStateMachine_Trigger_OnEnterActionNotFound(t *testing.T) {
	// Create a workflow definition with an OnEnter action that doesn't exist
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
				OnEnter: []string{
					"nonexistentAction",
				},
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, _, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent OnEnter action, got nil")
	}
}
