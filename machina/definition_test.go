package machina

import (
	"testing"
)

func TestState(t *testing.T) {
	// Create a state
	state := State{
		IsSideQuest: true,
		Name:        "testState",
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

	tests := []struct {
		name     string
		check    func() bool
		expected bool
		message  string
	}{
		{
			name:     "IsSideQuest",
			check:    func() bool { return state.IsSideQuest == true },
			expected: true,
			message:  "Expected IsSideQuest to be true",
		},
		{
			name:     "Name",
			check:    func() bool { return state.Name == "testState" },
			expected: true,
			message:  "Expected name to be 'testState'",
		},
		{
			name:     "OnEnterCount",
			check:    func() bool { return len(state.OnEnter) == 2 },
			expected: true,
			message:  "Expected 2 OnEnter actions",
		},
		{
			name:     "FirstOnEnterAction",
			check:    func() bool { return state.OnEnter[0] == "enterAction1" },
			expected: true,
			message:  "Expected first OnEnter action to be 'enterAction1'",
		},
		{
			name:     "SecondOnEnterAction",
			check:    func() bool { return state.OnEnter[1] == "enterAction2" },
			expected: true,
			message:  "Expected second OnEnter action to be 'enterAction2'",
		},
		{
			name:     "OnLeaveCount",
			check:    func() bool { return len(state.OnLeave) == 2 },
			expected: true,
			message:  "Expected 2 OnLeave actions",
		},
		{
			name:     "FirstOnLeaveAction",
			check:    func() bool { return state.OnLeave[0] == "leaveAction1" },
			expected: true,
			message:  "Expected first OnLeave action to be 'leaveAction1'",
		},
		{
			name:     "SecondOnLeaveAction",
			check:    func() bool { return state.OnLeave[1] == "leaveAction2" },
			expected: true,
			message:  "Expected second OnLeave action to be 'leaveAction2'",
		},
		{
			name:     "TransitionsCount",
			check:    func() bool { return len(state.Transitions) == 2 },
			expected: true,
			message:  "Expected 2 transitions",
		},
		{
			name:     "FirstTransitionEvent",
			check:    func() bool { return state.Transitions[0].Event == "event1" },
			expected: true,
			message:  "Expected first transition event to be 'event1'",
		},
		{
			name:     "FirstTransitionTarget",
			check:    func() bool { return state.Transitions[0].Target == "nextState1" },
			expected: true,
			message:  "Expected first transition target to be 'nextState1'",
		},
		{
			name:     "FirstTransitionConditionsCount",
			check:    func() bool { return len(state.Transitions[0].Conditions) == 2 },
			expected: true,
			message:  "Expected 2 conditions in first transition",
		},
		{
			name:     "FirstTransitionFirstCondition",
			check:    func() bool { return state.Transitions[0].Conditions[0] == "condition1" },
			expected: true,
			message:  "Expected first condition in first transition to be 'condition1'",
		},
		{
			name:     "FirstTransitionSecondCondition",
			check:    func() bool { return state.Transitions[0].Conditions[1] == "condition2" },
			expected: true,
			message:  "Expected second condition in first transition to be 'condition2'",
		},
		{
			name:     "FirstTransitionActionsCount",
			check:    func() bool { return len(state.Transitions[0].Actions) == 2 },
			expected: true,
			message:  "Expected 2 actions in first transition",
		},
		{
			name:     "FirstTransitionFirstAction",
			check:    func() bool { return state.Transitions[0].Actions[0] == "action1" },
			expected: true,
			message:  "Expected first action in first transition to be 'action1'",
		},
		{
			name:     "FirstTransitionSecondAction",
			check:    func() bool { return state.Transitions[0].Actions[1] == "action2" },
			expected: true,
			message:  "Expected second action in first transition to be 'action2'",
		},
		{
			name:     "SecondTransitionEvent",
			check:    func() bool { return state.Transitions[1].Event == "event2" },
			expected: true,
			message:  "Expected second transition event to be 'event2'",
		},
		{
			name:     "SecondTransitionTarget",
			check:    func() bool { return state.Transitions[1].Target == "nextState2" },
			expected: true,
			message:  "Expected second transition target to be 'nextState2'",
		},
		{
			name:     "SecondTransitionConditionsCount",
			check:    func() bool { return len(state.Transitions[1].Conditions) == 1 },
			expected: true,
			message:  "Expected 1 condition in second transition",
		},
		{
			name:     "SecondTransitionCondition",
			check:    func() bool { return state.Transitions[1].Conditions[0] == "condition3" },
			expected: true,
			message:  "Expected condition in second transition to be 'condition3'",
		},
		{
			name:     "SecondTransitionActionsCount",
			check:    func() bool { return len(state.Transitions[1].Actions) == 1 },
			expected: true,
			message:  "Expected 1 action in second transition",
		},
		{
			name:     "SecondTransitionAction",
			check:    func() bool { return state.Transitions[1].Actions[0] == "action3" },
			expected: true,
			message:  "Expected action in second transition to be 'action3'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.check()
			if result != tt.expected {
				t.Errorf("%s, got %v", tt.message, result)
			}
		})
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

	tests := []struct {
		name     string
		check    func() bool
		expected bool
		message  string
	}{
		{
			name:     "Event",
			check:    func() bool { return transition.Event == "testEvent" },
			expected: true,
			message:  "Expected event to be 'testEvent'",
		},
		{
			name:     "Target",
			check:    func() bool { return transition.Target == "nextState" },
			expected: true,
			message:  "Expected target to be 'nextState'",
		},
		{
			name:     "ConditionsCount",
			check:    func() bool { return len(transition.Conditions) == 2 },
			expected: true,
			message:  "Expected 2 conditions",
		},
		{
			name:     "FirstCondition",
			check:    func() bool { return transition.Conditions[0] == "condition1" },
			expected: true,
			message:  "Expected first condition to be 'condition1'",
		},
		{
			name:     "SecondCondition",
			check:    func() bool { return transition.Conditions[1] == "condition2" },
			expected: true,
			message:  "Expected second condition to be 'condition2'",
		},
		{
			name:     "ActionsCount",
			check:    func() bool { return len(transition.Actions) == 2 },
			expected: true,
			message:  "Expected 2 actions",
		},
		{
			name:     "FirstAction",
			check:    func() bool { return transition.Actions[0] == "action1" },
			expected: true,
			message:  "Expected first action to be 'action1'",
		},
		{
			name:     "SecondAction",
			check:    func() bool { return transition.Actions[1] == "action2" },
			expected: true,
			message:  "Expected second action to be 'action2'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.check()
			if result != tt.expected {
				t.Errorf("%s, got %v", tt.message, result)
			}
		})
	}
}

func TestWorkflowDefinition(t *testing.T) {
	// Create a workflow definition
	workflow := WorkflowDefinition{
		States: map[string]State{
			"state1": {
				IsSideQuest: false,
				Name:        "state1",
			},
			"state2": {
				IsSideQuest: true,
				Name:        "state2",
			},
		},
	}

	tests := []struct {
		name     string
		check    func() bool
		expected bool
		message  string
	}{
		{
			name:     "StatesCount",
			check:    func() bool { return len(workflow.States) == 2 },
			expected: true,
			message:  "Expected 2 states",
		},
		{
			name:     "FirstStateName",
			check:    func() bool { return workflow.States["state1"].Name == "state1" },
			expected: true,
			message:  "Expected first state name to be 'state1'",
		},
		{
			name:     "FirstStateIsSideQuest",
			check:    func() bool { return workflow.States["state1"].IsSideQuest == false },
			expected: true,
			message:  "Expected first state IsSideQuest to be false",
		},
		{
			name:     "SecondStateName",
			check:    func() bool { return workflow.States["state2"].Name == "state2" },
			expected: true,
			message:  "Expected second state name to be 'state2'",
		},
		{
			name:     "SecondStateIsSideQuest",
			check:    func() bool { return workflow.States["state2"].IsSideQuest == true },
			expected: true,
			message:  "Expected second state IsSideQuest to be true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.check()
			if result != tt.expected {
				t.Errorf("%s, got %v", tt.message, result)
			}
		})
	}
}