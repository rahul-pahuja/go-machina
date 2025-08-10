package machina

import (
	"testing"
)

func TestWorkflowDefinition_Validate(t *testing.T) {
	tests := []struct {
		name        string
		definition  *WorkflowDefinition
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidWorkflow",
			definition: &WorkflowDefinition{
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
			},
			expectError: false,
		},
		{
			name: "ValidWorkflowWithInitialState",
			definition: &WorkflowDefinition{
				InitialState: "start",
				States: map[string]State{
					"start": {
						Name: "start",
					},
					"end": {
						Name: "end",
					},
				},
			},
			expectError: false,
		},
		{
			name: "EmptyWorkflow",
			definition: &WorkflowDefinition{
				States: map[string]State{},
			},
			expectError: true,
			errorMsg:    "workflow must have at least one state",
		},
		{
			name: "MismatchedStateKeyAndName",
			definition: &WorkflowDefinition{
				States: map[string]State{
					"start": {
						Name: "different",
					},
				},
			},
			expectError: true,
			errorMsg:    "state key start does not match state name different",
		},
		{
			name: "InvalidInitialState",
			definition: &WorkflowDefinition{
				InitialState: "nonexistent",
				States: map[string]State{
					"start": {
						Name: "start",
					},
				},
			},
			expectError: true,
			errorMsg:    "initialState nonexistent not found in states",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.definition.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestState_Validate(t *testing.T) {
	tests := []struct {
		name        string
		state       *State
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidState",
			state: &State{
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			expectError: false,
		},
		{
			name: "ValidStateWithSideQuest",
			state: &State{
				Name:        "sideQuest",
				IsSideQuest: true,
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			expectError: false,
		},
		{
			name: "StateWithNoName",
			state: &State{
				Name: "",
				Transitions: []Transition{
					{
						Event:  "proceed",
						Target: "end",
					},
				},
			},
			expectError: true,
			errorMsg:    "state must have a name",
		},
		{
			name: "StateWithInvalidTransition",
			state: &State{
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "",
						Target: "end",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid transition for event : transition must have an event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.state.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestTransition_Validate(t *testing.T) {
	tests := []struct {
		name        string
		transition  *Transition
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidTransition",
			transition: &Transition{
				Event:  "proceed",
				Target: "end",
			},
			expectError: false,
		},
		{
			name: "TransitionWithNoEvent",
			transition: &Transition{
				Event:  "",
				Target: "end",
			},
			expectError: true,
			errorMsg:    "transition must have an event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transition.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestWorkflowDefinition_InitialState(t *testing.T) {
	tests := []struct {
		name        string
		definition  *WorkflowDefinition
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidInitialState",
			definition: &WorkflowDefinition{
				InitialState: "start",
				States: map[string]State{
					"start": {
						Name: "start",
					},
					"end": {
						Name: "end",
					},
				},
			},
			expectError: false,
		},
		{
			name: "InvalidInitialState",
			definition: &WorkflowDefinition{
				InitialState: "nonexistent",
				States: map[string]State{
					"start": {
						Name: "start",
					},
				},
			},
			expectError: true,
			errorMsg:    "initialState nonexistent not found in states",
		},
		{
			name: "EmptyInitialState",
			definition: &WorkflowDefinition{
				InitialState: "",
				States: map[string]State{
					"start": {
						Name: "start",
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.definition.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}