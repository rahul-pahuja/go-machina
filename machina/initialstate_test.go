package machina

import (
	"testing"
)

func TestWorkflowDefinition_InitialState(t *testing.T) {
	// Test with valid initial state
	definition := &WorkflowDefinition{
		InitialState: "start",
		States: map[string]State{
			"start": {
				Name: "start",
			},
			"end": {
				Name: "end",
			},
		},
	}

	if err := definition.Validate(); err != nil {
		t.Errorf("Expected valid definition, got error: %v", err)
	}

	if definition.InitialState != "start" {
		t.Errorf("Expected InitialState to be 'start', got %s", definition.InitialState)
	}

	// Test with invalid initial state
	definition2 := &WorkflowDefinition{
		InitialState: "nonexistent",
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	if err := definition2.Validate(); err == nil {
		t.Error("Expected validation error for invalid initial state, got nil")
	}

	// Test with empty initial state (should be valid)
	definition3 := &WorkflowDefinition{
		InitialState: "",
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	if err := definition3.Validate(); err != nil {
		t.Errorf("Expected valid definition with empty initial state, got error: %v", err)
	}
}