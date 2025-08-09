package machina

import (
	"testing"
)

func TestStateMachine_GetAutoEventForTransition(t *testing.T) {
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

	// Get auto event for transition
	autoEvent, err := fsm.GetAutoEventForTransition("start", "proceed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if autoEvent != "completed" {
		t.Errorf("Expected auto event to be 'completed', got %s", autoEvent)
	}
}

func TestStateMachine_GetAutoEventForTransition_NoAutoEvent(t *testing.T) {
	// Create a workflow definition with a transition that has no auto event
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

	// Get auto event for transition
	autoEvent, err := fsm.GetAutoEventForTransition("start", "proceed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if autoEvent != "" {
		t.Errorf("Expected no auto event, got %s", autoEvent)
	}
}
