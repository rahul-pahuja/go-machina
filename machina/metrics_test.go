package machina

import (
	"context"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestMetrics(t *testing.T) {
	// Create a test registry
	reg := prometheus.NewRegistry()

	// Create a simple workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "next",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry with a simple action
	registry := NewRegistry()
	registry.RegisterAction("testAction", func(ctx context.Context, data map[string]any) (map[string]any, error) {
		return map[string]any{"result": "success"}, nil
	})

	// Create a logger
	logger := slog.Default()

	// Create the state machine with metrics
	sm := NewStateMachine(definition, registry, logger, WithMetrics(reg), WithTracer(noop.NewTracerProvider().Tracer("test")))

	// Perform a transition
	_, err := sm.Trigger(context.Background(), "start", "next", map[string]any{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that we can gather metrics without error
	_, err = reg.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}
}

func TestMetricsError(t *testing.T) {
	// Create a test registry
	reg := prometheus.NewRegistry()

	// Create a simple workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "next",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry (no need for special conditions)
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine with metrics
	sm := NewStateMachine(definition, registry, logger, WithMetrics(reg), WithTracer(noop.NewTracerProvider().Tracer("test")))

	// Try to perform a transition with a non-existent event
	_, err := sm.Trigger(context.Background(), "start", "nonexistent", map[string]any{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that we can gather metrics without error
	_, err = reg.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}
}

func TestMetricsAutoTransition(t *testing.T) {
	// Create a test registry
	reg := prometheus.NewRegistry()

	// Create a workflow definition with an auto transition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:     "next",
						Target:    "end",
						AutoEvent: "auto",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine with metrics
	sm := NewStateMachine(definition, registry, logger, WithMetrics(reg), WithTracer(noop.NewTracerProvider().Tracer("test")))

	// Perform a transition
	_, err := sm.Trigger(context.Background(), "start", "next", map[string]any{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that we can gather metrics without error
	_, err = reg.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}
}

func TestGetAutoEventForTransition(t *testing.T) {
	// Create a workflow definition with an auto transition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:     "next",
						Target:    "end",
						AutoEvent: "auto",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine
	sm := NewStateMachine(definition, registry, logger)

	// Get the auto event
	autoEvent, err := sm.GetAutoEventForTransition("start", "next")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if autoEvent != "auto" {
		t.Errorf("Expected auto event 'auto', got '%s'", autoEvent)
	}
}

func TestGetAutoEventForTransition_NoAutoEvent(t *testing.T) {
	// Create a workflow definition without an auto transition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "next",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine
	sm := NewStateMachine(definition, registry, logger)

	// Get the auto event
	autoEvent, err := sm.GetAutoEventForTransition("start", "next")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if autoEvent != "" {
		t.Errorf("Expected no auto event, got '%s'", autoEvent)
	}
}

func TestGetAutoEventForTransition_StateNotFound(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "next",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine
	sm := NewStateMachine(definition, registry, logger)

	// Try to get auto event for non-existent state
	_, err := sm.GetAutoEventForTransition("nonexistent", "next")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetAutoEventForTransition_TransitionNotFound(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
				Transitions: []Transition{
					{
						Event:  "next",
						Target: "end",
					},
				},
			},
			"end": {
				Name: "end",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine
	sm := NewStateMachine(definition, registry, logger)

	// Try to get auto event for non-existent transition
	_, err := sm.GetAutoEventForTransition("start", "nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestNewMetrics(t *testing.T) {
	// Create a test registry
	reg := prometheus.NewRegistry()

	// Create metrics
	metrics := NewMetrics(reg)

	// Verify that metrics were created
	if metrics.TransitionsTotal == nil {
		t.Error("TransitionsTotal metric not created")
	}

	if metrics.TransitionErrors == nil {
		t.Error("TransitionErrors metric not created")
	}

	if metrics.TransitionDuration == nil {
		t.Error("TransitionDuration metric not created")
	}

	if metrics.AutoTransitionsTotal == nil {
		t.Error("AutoTransitionsTotal metric not created")
	}
}
