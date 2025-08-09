package machina

import (
	"context"
	"testing"
)

func TestStateMachine_Trigger_WithAutoEvent(t *testing.T) {
	// Create a workflow definition with a transition that has an auto event
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:     "proceed",
						Target:    "end",
						AutoEvent: "completed",
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

	if result.AutoEvent != "completed" {
		t.Errorf("Expected auto event to be 'completed', got %s", result.AutoEvent)
	}

	if result.PersistenceData == nil {
		t.Error("Expected persistence data, got nil")
	}
}