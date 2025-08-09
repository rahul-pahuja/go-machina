package machina

import (
	"testing"
)

func TestState(t *testing.T) {
	// Create a state
	state := State{
		Name: "testState",
		OnEnter: []string{
			"enterAction1",
			"enterAction2",
		},
		OnLeave: []string{
			"leaveAction1",
			"leaveAction2",
		},
		Transitions: []Transition{
			{
				Event:  "event1",
				Target: "nextState1",
				Conditions: []string{
					"condition1",
					"condition2",
				},
				Actions: []string{
					"action1",
					"action2",
				},
			},
			{
				Event:  "event2",
				Target: "nextState2",
				Conditions: []string{
					"condition3",
				},
				Actions: []string{
					"action3",
				},
			},
		},
	}

	// Verify the state
	if state.Name != "testState" {
		t.Errorf("Expected name to be 'testState', got %s", state.Name)
	}

	if len(state.OnEnter) != 2 {
		t.Errorf("Expected 2 OnEnter actions, got %d", len(state.OnEnter))
	}

	if state.OnEnter[0] != "enterAction1" {
		t.Errorf("Expected first OnEnter action to be 'enterAction1', got %s", state.OnEnter[0])
	}

	if state.OnEnter[1] != "enterAction2" {
		t.Errorf("Expected second OnEnter action to be 'enterAction2', got %s", state.OnEnter[1])
	}

	if len(state.OnLeave) != 2 {
		t.Errorf("Expected 2 OnLeave actions, got %d", len(state.OnLeave))
	}

	if state.OnLeave[0] != "leaveAction1" {
		t.Errorf("Expected first OnLeave action to be 'leaveAction1', got %s", state.OnLeave[0])
	}

	if state.OnLeave[1] != "leaveAction2" {
		t.Errorf("Expected second OnLeave action to be 'leaveAction2', got %s", state.OnLeave[1])
	}

	if len(state.Transitions) != 2 {
		t.Errorf("Expected 2 transitions, got %d", len(state.Transitions))
	}

	if state.Transitions[0].Event != "event1" {
		t.Errorf("Expected first transition event to be 'event1', got %s", state.Transitions[0].Event)
	}

	if state.Transitions[0].Target != "nextState1" {
		t.Errorf("Expected first transition target to be 'nextState1', got %s", state.Transitions[0].Target)
	}

	if len(state.Transitions[0].Conditions) != 2 {
		t.Errorf("Expected 2 conditions in first transition, got %d", len(state.Transitions[0].Conditions))
	}

	if state.Transitions[0].Conditions[0] != "condition1" {
		t.Errorf("Expected first condition in first transition to be 'condition1', got %s", state.Transitions[0].Conditions[0])
	}

	if state.Transitions[0].Conditions[1] != "condition2" {
		t.Errorf("Expected second condition in first transition to be 'condition2', got %s", state.Transitions[0].Conditions[1])
	}

	if len(state.Transitions[0].Actions) != 2 {
		t.Errorf("Expected 2 actions in first transition, got %d", len(state.Transitions[0].Actions))
	}

	if state.Transitions[0].Actions[0] != "action1" {
		t.Errorf("Expected first action in first transition to be 'action1', got %s", state.Transitions[0].Actions[0])
	}

	if state.Transitions[0].Actions[1] != "action2" {
		t.Errorf("Expected second action in first transition to be 'action2', got %s", state.Transitions[0].Actions[1])
	}

	if state.Transitions[1].Event != "event2" {
		t.Errorf("Expected second transition event to be 'event2', got %s", state.Transitions[1].Event)
	}

	if state.Transitions[1].Target != "nextState2" {
		t.Errorf("Expected second transition target to be 'nextState2', got %s", state.Transitions[1].Target)
	}

	if len(state.Transitions[1].Conditions) != 1 {
		t.Errorf("Expected 1 condition in second transition, got %d", len(state.Transitions[1].Conditions))
	}

	if state.Transitions[1].Conditions[0] != "condition3" {
		t.Errorf("Expected condition in second transition to be 'condition3', got %s", state.Transitions[1].Conditions[0])
	}

	if len(state.Transitions[1].Actions) != 1 {
		t.Errorf("Expected 1 action in second transition, got %d", len(state.Transitions[1].Actions))
	}

	if state.Transitions[1].Actions[0] != "action3" {
		t.Errorf("Expected action in second transition to be 'action3', got %s", state.Transitions[1].Actions[0])
	}
}

func TestTransition(t *testing.T) {
	// Create a transition
	transition := Transition{
		Event:  "testEvent",
		Target: "nextState",
		Conditions: []string{
			"condition1",
			"condition2",
		},
		Actions: []string{
			"action1",
			"action2",
		},
	}

	// Verify the transition
	if transition.Event != "testEvent" {
		t.Errorf("Expected event to be 'testEvent', got %s", transition.Event)
	}

	if transition.Target != "nextState" {
		t.Errorf("Expected target to be 'nextState', got %s", transition.Target)
	}

	if len(transition.Conditions) != 2 {
		t.Errorf("Expected 2 conditions, got %d", len(transition.Conditions))
	}

	if transition.Conditions[0] != "condition1" {
		t.Errorf("Expected first condition to be 'condition1', got %s", transition.Conditions[0])
	}

	if transition.Conditions[1] != "condition2" {
		t.Errorf("Expected second condition to be 'condition2', got %s", transition.Conditions[1])
	}

	if len(transition.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(transition.Actions))
	}

	if transition.Actions[0] != "action1" {
		t.Errorf("Expected first action to be 'action1', got %s", transition.Actions[0])
	}

	if transition.Actions[1] != "action2" {
		t.Errorf("Expected second action to be 'action2', got %s", transition.Actions[1])
	}
}

func TestWorkflowDefinition(t *testing.T) {
	// Create a workflow definition
	workflow := WorkflowDefinition{
		States: map[string]State{
			"state1": {
				Name: "state1",
			},
			"state2": {
				Name: "state2",
			},
		},
	}

	// Verify the workflow definition
	if len(workflow.States) != 2 {
		t.Errorf("Expected 2 states, got %d", len(workflow.States))
	}

	if workflow.States["state1"].Name != "state1" {
		t.Errorf("Expected first state name to be 'state1', got %s", workflow.States["state1"].Name)
	}

	if workflow.States["state2"].Name != "state2" {
		t.Errorf("Expected second state name to be 'state2', got %s", workflow.States["state2"].Name)
	}
}
