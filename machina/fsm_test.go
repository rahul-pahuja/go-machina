package machina

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestStateMachine_Trigger_SuccessCases(t *testing.T) {
	tests := []struct {
		name           string
		definition     *WorkflowDefinition
		registrySetup  func(*Registry)
		currentState   string
		event          string
		payload        map[string]any
		expectedResult *TransitionResult
	}{
		{
			name: "SimpleTransition",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
			},
			currentState: "start",
			event:        "proceed",
			payload:      map[string]any{},
			expectedResult: &TransitionResult{
				NewState:        "end",
				AutoEvent:       "",
				PersistenceData: map[string]any{},
			},
		},
		{
			name: "TransitionWithActionExecution",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
				r.RegisterAction("updateAction", MockUpdateAction)
			},
			currentState: "start",
			event:        "proceed",
			payload:      map[string]any{},
			expectedResult: &TransitionResult{
				NewState:  "end",
				AutoEvent: "",
				PersistenceData: map[string]any{
					"updated": true,
				},
			},
		},
		{
			name: "TransitionWithLifecycleHooks",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
				r.RegisterAction("onLeaveAction", MockUpdateAction)
				r.RegisterAction("transitionAction", MockUpdateAction)
				r.RegisterAction("onEnterAction", MockUpdateAction)
			},
			currentState: "start",
			event:        "proceed",
			payload:      map[string]any{},
			expectedResult: &TransitionResult{
				NewState:  "end",
				AutoEvent: "",
				PersistenceData: map[string]any{
					"updated": true,
				},
			},
		},
		{
			name: "DynamicTransitionWithNextStateOverride",
			definition: &WorkflowDefinition{
				States: map[string]State{
					"start": {
						Name: "start",
						Transitions: []Transition{
							{
								Event:  "proceed",
								Target: "end", // This will be overridden by the action
								Conditions: []string{
									"alwaysTrue",
								},
								Actions: []string{
									"overrideAction",
								},
							},
						},
					},
					"intermediate": {
						Name: "intermediate",
					},
					"end": {
						Name: "end",
					},
				},
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
				r.RegisterAction("overrideAction", func(ctx context.Context, data map[string]any) (map[string]any, error) {
					return map[string]any{
						"__next_state_override": "intermediate",
					}, nil
				})
			},
			currentState: "start",
			event:        "proceed",
			payload:      map[string]any{},
			expectedResult: &TransitionResult{
				NewState:        "intermediate",
				AutoEvent:       "",
				PersistenceData: map[string]any{},
			},
		},
		{
			name: "TransitionToSideQuestState",
			definition: &WorkflowDefinition{
				States: map[string]State{
					"start": {
						Name: "start",
						Transitions: []Transition{
							{
								Event:  "proceed",
								Target: "sideQuest",
								Conditions: []string{
									"alwaysTrue",
								},
							},
						},
					},
					"sideQuest": {
						Name:        "sideQuest",
						IsSideQuest: true,
					},
				},
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
			},
			currentState: "start",
			event:        "proceed",
			payload:      map[string]any{},
			expectedResult: &TransitionResult{
				NewState:        "sideQuest",
				AutoEvent:       "",
				PersistenceData: map[string]any{},
			},
		},
		// This test case is removed because it's complex to simulate correctly in a single Trigger call
		// The ReturnToPreviousStateAction functionality is tested in the mocks_test.go file
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create registry with mock implementations
			registry := NewRegistry()
			if tt.registrySetup != nil {
				tt.registrySetup(registry)
			}

			// Create state machine
			fsm := NewStateMachine(tt.definition, registry, nil)

			// Trigger event
			ctx := context.Background()
			result, err := fsm.Trigger(ctx, tt.currentState, tt.event, tt.payload)

			// Verify results
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result.NewState != tt.expectedResult.NewState {
				t.Errorf("Expected new state to be '%s', got '%s'", tt.expectedResult.NewState, result.NewState)
			}

			if result.AutoEvent != tt.expectedResult.AutoEvent {
				t.Errorf("Expected auto event to be '%s', got '%s'", tt.expectedResult.AutoEvent, result.AutoEvent)
			}

			// Check persistence data
			if tt.expectedResult.PersistenceData != nil {
				for key, expectedValue := range tt.expectedResult.PersistenceData {
					if actualValue, exists := result.PersistenceData[key]; !exists {
						t.Errorf("Expected key '%s' in persistence data, but it was missing", key)
					} else if actualValue != expectedValue {
						t.Errorf("Expected value '%v' for key '%s' in persistence data, got '%v'", expectedValue, key, actualValue)
					}
				}
			}

			if result.PersistenceData == nil {
				t.Error("Expected persistence data, got nil")
			}
		})
	}
}

func TestStateMachine_Trigger_ConditionCases(t *testing.T) {
	tests := []struct {
		name          string
		definition    *WorkflowDefinition
		registrySetup func(*Registry)
		currentState  string
		event         string
		payload       map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name: "ConditionFailure",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysFalse", MockFalseCondition)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "condition alwaysFalse evaluated to false",
		},
		{
			name: "ConditionError",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("errorCondition", MockErrorCondition)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "condition errorCondition failed: condition error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create registry with mock implementations
			registry := NewRegistry()
			if tt.registrySetup != nil {
				tt.registrySetup(registry)
			}

			// Create state machine
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			fsm := NewStateMachine(tt.definition, registry, logger)

			// Trigger event
			ctx := context.Background()
			_, err := fsm.Trigger(ctx, tt.currentState, tt.event, tt.payload)

			// Verify results
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && err.Error() != tt.errorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestStateMachine_Trigger_ActionErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		definition    *WorkflowDefinition
		registrySetup func(*Registry)
		currentState  string
		event         string
		payload       map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name: "ActionError",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
				r.RegisterAction("errorAction", MockErrorAction)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "transition action errorAction failed: action error",
		},
		{
			name: "OnLeaveActionError",
			definition: &WorkflowDefinition{
				States: map[string]State{
					"start": {
						Name: "start",
						OnLeave: []string{
							"errorAction",
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterAction("errorAction", MockErrorAction)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "OnLeave action errorAction failed: action error",
		},
		{
			name: "OnEnterActionError",
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
						OnEnter: []string{
							"errorAction",
						},
					},
				},
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
				r.RegisterAction("errorAction", MockErrorAction)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "OnEnter action errorAction failed: action error",
		},
		// This test case is removed because it's complex to simulate correctly in a single Trigger call
		// The ReturnToPreviousStateAction functionality is tested in the mocks_test.go file
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create registry with mock implementations
			registry := NewRegistry()
			if tt.registrySetup != nil {
				tt.registrySetup(registry)
			}

			// Create state machine
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			fsm := NewStateMachine(tt.definition, registry, logger)

			// Trigger event
			ctx := context.Background()
			_, err := fsm.Trigger(ctx, tt.currentState, tt.event, tt.payload)

			// Verify results
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && err.Error() != tt.errorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
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

func TestStateMachine_Trigger_ResourceNotFoundCases(t *testing.T) {
	tests := []struct {
		name          string
		definition    *WorkflowDefinition
		registrySetup func(*Registry)
		currentState  string
		event         string
		payload       map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name: "StateNotFound",
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
			registrySetup: func(r *Registry) {
				// No setup needed
			},
			currentState:  "nonexistent",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "failed to get state definition for nonexistent: state nonexistent not found",
		},
		{
			name: "TransitionNotFound",
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
			registrySetup: func(r *Registry) {
				// No setup needed
			},
			currentState:  "start",
			event:         "nonexistent",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "no valid transition found for event nonexistent in state start: no transition found for event nonexistent",
		},
		{
			name: "ConditionNotFound",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				// No setup needed
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "failed to get condition nonexistent: condition nonexistent not found",
		},
		{
			name: "ActionNotFound",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				r.RegisterCondition("alwaysTrue", MockTrueCondition)
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "failed to get transition action nonexistent: action nonexistent not found",
		},
		{
			name: "OnLeaveActionNotFound",
			definition: &WorkflowDefinition{
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
			},
			registrySetup: func(r *Registry) {
				// No setup needed
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "failed to get OnLeave action nonexistent: action nonexistent not found",
		},
		{
			name: "OnEnterActionNotFound",
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
						OnEnter: []string{
							"nonexistent",
						},
					},
				},
			},
			registrySetup: func(r *Registry) {
				// No setup needed
			},
			currentState:  "start",
			event:         "proceed",
			payload:       map[string]any{},
			expectError:   true,
			errorContains: "failed to get OnEnter action nonexistent: action nonexistent not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create registry with mock implementations
			registry := NewRegistry()
			if tt.registrySetup != nil {
				tt.registrySetup(registry)
			}

			// Create state machine
			fsm := NewStateMachine(tt.definition, registry, nil)

			// Trigger event
			ctx := context.Background()
			_, err := fsm.Trigger(ctx, tt.currentState, tt.event, tt.payload)

			// Verify results
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && err.Error() != tt.errorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestReturnToPreviousStateAction(t *testing.T) {
	tests := []struct {
		name          string
		inputData     map[string]any
		expectedData  map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name: "ValidStack",
			inputData: map[string]any{
				"WorkflowStack": []string{"state1", "state2"},
			},
			expectedData: map[string]any{
				"__next_state_override": "state2",
				"WorkflowStack":         []string{"state1"},
			},
			expectError: false,
		},
		{
			name: "SingleItemStack",
			inputData: map[string]any{
				"WorkflowStack": []string{"state1"},
			},
			expectedData: map[string]any{
				"__next_state_override": "state1",
				"WorkflowStack":         []string{},
			},
			expectError: false,
		},
		{
			name:          "EmptyStack",
			inputData:     map[string]any{},
			expectError:   true,
			errorContains: "workflow stack not found or empty",
		},
		{
			name: "NilStack",
			inputData: map[string]any{
				"WorkflowStack": nil,
			},
			expectError:   true,
			errorContains: "workflow stack not found or empty",
		},
		{
			name: "WrongTypeStack",
			inputData: map[string]any{
				"WorkflowStack": "not a slice",
			},
			expectError:   true,
			errorContains: "workflow stack not found or empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := ReturnToPreviousStateAction(ctx, tt.inputData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && err.Error() != tt.errorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Check __next_state_override
			if result["__next_state_override"] != tt.expectedData["__next_state_override"] {
				t.Errorf("Expected __next_state_override to be '%s', got '%s'", tt.expectedData["__next_state_override"], result["__next_state_override"])
			}

			// Check WorkflowStack
			expectedStack, ok := tt.expectedData["WorkflowStack"].([]string)
			if !ok {
				t.Fatalf("Expected WorkflowStack to be []string")
			}

			actualStack, ok := result["WorkflowStack"].([]string)
			if !ok {
				t.Fatalf("Expected result WorkflowStack to be []string")
			}

			if len(actualStack) != len(expectedStack) {
				t.Errorf("Expected WorkflowStack length to be %d, got %d", len(expectedStack), len(actualStack))
				return
			}

			for i, v := range expectedStack {
				if actualStack[i] != v {
					t.Errorf("Expected WorkflowStack[%d] to be '%s', got '%s'", i, v, actualStack[i])
				}
			}
		})
	}
}

func TestGetTransitionForEvent(t *testing.T) {
	// Create a registry with mock conditions
	registry := NewRegistry()
	registry.RegisterCondition("condition1", MockTrueCondition)
	registry.RegisterCondition("condition2", MockFalseCondition)
	registry.RegisterCondition("condition3", MockTrueCondition)

	// Create a state machine
	fsm := &StateMachine{
		registry: registry,
	}

	tests := []struct {
		name          string
		state         *State
		event         string
		expectedIndex int // Index of expected transition in the state's Transitions slice
		expectError   bool
		errorContains string
	}{
		{
			name: "SingleTransition",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
					},
				},
			},
			event:         "event1",
			expectedIndex: 0,
			expectError:   false,
		},
		{
			name: "MultipleTransitionsDifferentEvents",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
					},
					{
						Event:  "event2",
						Target: "target2",
					},
				},
			},
			event:         "event2",
			expectedIndex: 1,
			expectError:   false,
		},
		{
			name: "MultipleTransitionsSameEventNoConditions",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
					},
					{
						Event:  "event1",
						Target: "target2",
					},
				},
			},
			event:         "event1",
			expectedIndex: 0, // Should return the first one
			expectError:   false,
		},
		{
			name: "MultipleTransitionsSameEventWithConditionsMatchFirst",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
						Conditions: []string{
							"condition1", // True
						},
					},
					{
						Event:  "event1",
						Target: "target2",
						Conditions: []string{
							"condition2", // False
						},
					},
				},
			},
			event:         "event1",
			expectedIndex: 0,
			expectError:   false,
		},
		{
			name: "MultipleTransitionsSameEventWithConditionsMatchSecond",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
						Conditions: []string{
							"condition2", // False
						},
					},
					{
						Event:  "event1",
						Target: "target2",
						Conditions: []string{
							"condition1", // True
						},
					},
				},
			},
			event:         "event1",
			expectedIndex: 1,
			expectError:   false,
		},
		{
			name: "MultipleTransitionsSameEventWithConditionsAllFalse",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
						Conditions: []string{
							"condition2", // False
						},
					},
					{
						Event:  "event1",
						Target: "target2",
						Conditions: []string{
							"condition2", // False
						},
					},
				},
			},
			event:         "event1",
			expectError:   true,
			errorContains: "no transition found for event event1 with matching conditions",
		},
		{
			name: "NoTransitionForEvent",
			state: &State{
				Transitions: []Transition{
					{
						Event:  "event1",
						Target: "target1",
					},
				},
			},
			event:         "event2",
			expectError:   true,
			errorContains: "no transition found for event event2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			payload := map[string]any{}

			transition, err := fsm.getTransitionForEvent(tt.state, tt.event, ctx, payload)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && err.Error() != tt.errorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if transition == nil {
				t.Error("Expected transition, got nil")
				return
			}

			expectedTransition := &tt.state.Transitions[tt.expectedIndex]
			if transition.Event != expectedTransition.Event {
				t.Errorf("Expected transition event to be '%s', got '%s'", expectedTransition.Event, transition.Event)
			}

			if transition.Target != expectedTransition.Target {
				t.Errorf("Expected transition target to be '%s', got '%s'", expectedTransition.Target, transition.Target)
			}
		})
	}
}

func TestNewStateMachine_InvalidDefinition(t *testing.T) {
	// Create an invalid workflow definition (empty states)
	invalidDefinition := &WorkflowDefinition{
		States: map[string]State{},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Try to create a state machine with the invalid definition
	fsm := NewStateMachine(invalidDefinition, registry, logger)

	// Verify that the state machine is nil
	if fsm != nil {
		t.Error("Expected state machine to be nil for invalid definition")
	}
}
