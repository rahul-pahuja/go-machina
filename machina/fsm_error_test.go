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
	_, err := fsm.Trigger(context.Background(), "nonexistent", "proceed", map[string]any{})

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
	_, err := fsm.Trigger(context.Background(), "start", "nonexistent", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent transition, got nil")
	}
}

func TestStateMachine_Trigger_ConditionNotFound(t *testing.T) {
	// Create a workflow definition with a non-existent condition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Conditions: []string{
							"nonexistent",
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
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent condition, got nil")
	}
}

func TestStateMachine_Trigger_ActionNotFound(t *testing.T) {
	// Create a workflow definition with a non-existent action
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
						Actions: []string{
							"nonexistent",
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
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent action, got nil")
	}
}

func TestStateMachine_Trigger_OnLeaveActionNotFound(t *testing.T) {
	// Create a workflow definition with a non-existent OnLeave action
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				OnLeave: []string{
					"nonexistent",
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
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent OnLeave action, got nil")
	}
}

func TestStateMachine_Trigger_OnEnterActionNotFound(t *testing.T) {
	// Create a workflow definition with a non-existent OnEnter action
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
					"nonexistent",
				},
			},
		},
	}

	// Create registry
	registry := NewRegistry()

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Try to trigger event
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})

	// Verify results
	if err == nil {
		t.Error("Expected error for non-existent OnEnter action, got nil")
	}
}

func TestStateMachine_Trigger_GuardConditionFailure(t *testing.T) {
	// Create a workflow definition
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

	// Try to trigger event with failing guard
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{}, MockFalseCondition)

	// Verify results
	if err == nil {
		t.Error("Expected error for failing guard condition, got nil")
	}
}

func TestStateMachine_Trigger_GuardConditionError(t *testing.T) {
	// Create a workflow definition
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

	// Try to trigger event with erroring guard
	_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{}, MockErrorCondition)

	// Verify results
	if err == nil {
		t.Error("Expected error for erroring guard condition, got nil")
	}
}
