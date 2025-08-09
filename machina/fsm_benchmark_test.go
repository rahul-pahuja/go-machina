package machina

import (
	"context"
	"testing"
)

func BenchmarkStateMachine_Trigger(b *testing.B) {
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
						Actions: []string{
							"noOpAction",
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
	registry.RegisterAction("noOpAction", MockNoOpAction)

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStateMachine_Trigger_WithGuards(b *testing.B) {
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
						Actions: []string{
							"noOpAction",
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
	registry.RegisterAction("noOpAction", MockNoOpAction)

	// Create state machine
	fsm := NewStateMachine(definition, registry, nil)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fsm.Trigger(context.Background(), "start", "proceed", map[string]any{}, MockGuardCondition)
		if err != nil {
			b.Fatal(err)
		}
	}
}
