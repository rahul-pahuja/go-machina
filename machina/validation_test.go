package machina

import (
	"testing"
)

func TestWorkflowDefinition_Validate(t *testing.T) {
	// Test valid workflow definition
	validWorkflow := &WorkflowDefinition{
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

	if err := validWorkflow.Validate(); err != nil {
		t.Errorf("Expected valid workflow, got error: %v", err)
	}

	// Test workflow with no states
	emptyWorkflow := &WorkflowDefinition{
		States: map[string]State{},
	}

	if err := emptyWorkflow.Validate(); err == nil {
		t.Error("Expected error for empty workflow, got nil")
	}

	// Test workflow with mismatched state key and name
	mismatchedWorkflow := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "different",
			},
		},
	}

	if err := mismatchedWorkflow.Validate(); err == nil {
		t.Error("Expected error for mismatched state key and name, got nil")
	}
}

func TestState_Validate(t *testing.T) {
	// Test valid state
	validState := &State{
		Name: "start",
		Transitions: []Transition{
			{
				Event:  "proceed",
				Target: "end",
			},
		},
	}

	if err := validState.Validate(); err != nil {
		t.Errorf("Expected valid state, got error: %v", err)
	}

	// Test state with no name
	invalidState := &State{
		Name: "",
		Transitions: []Transition{
			{
				Event:  "proceed",
				Target: "end",
			},
		},
	}

	if err := invalidState.Validate(); err == nil {
		t.Error("Expected error for invalid state, got nil")
	}

	// Test state with invalid transition (no event)
	invalidTransitionState := &State{
		Name: "start",
		Transitions: []Transition{
			{
				Event:  "",
				Target: "end",
			},
		},
	}

	if err := invalidTransitionState.Validate(); err == nil {
		t.Error("Expected error for invalid transition, got nil")
	}

	// Test state with invalid transition (no target)
	invalidTransitionState2 := &State{
		Name: "start",
		Transitions: []Transition{
			{
				Event:  "proceed",
				Target: "",
			},
		},
	}

	if err := invalidTransitionState2.Validate(); err == nil {
		t.Error("Expected error for invalid transition, got nil")
	}
}

func TestTransition_Validate(t *testing.T) {
	// Test valid transition
	validTransition := &Transition{
		Event:  "proceed",
		Target: "end",
	}

	if err := validTransition.Validate(); err != nil {
		t.Errorf("Expected valid transition, got error: %v", err)
	}

	// Test transition with no event
	invalidTransition := &Transition{
		Event:  "",
		Target: "end",
	}

	if err := invalidTransition.Validate(); err == nil {
		t.Error("Expected error for transition with no event, got nil")
	}

	// Test transition with no target
	invalidTransition2 := &Transition{
		Event:  "proceed",
		Target: "",
	}

	if err := invalidTransition2.Validate(); err == nil {
		t.Error("Expected error for transition with no target, got nil")
	}
}
