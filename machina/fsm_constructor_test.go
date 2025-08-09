package machina

import (
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNewStateMachine(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine
	sm := NewStateMachine(definition, registry, logger)

	// Verify that the state machine was created
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}

	// Verify that the fields were set correctly
	if sm.definition != definition {
		t.Error("Definition not set correctly")
	}

	if sm.registry != registry {
		t.Error("Registry not set correctly")
	}

	if sm.logger != logger {
		t.Error("Logger not set correctly")
	}
}

func TestNewStateMachine_WithMetrics(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create a prometheus registry
	promRegistry := prometheus.NewRegistry()

	// Create the state machine with metrics
	sm := NewStateMachine(definition, registry, logger, WithMetrics(promRegistry))

	// Verify that the state machine was created
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}

	// Verify that metrics were set
	if sm.metrics == nil {
		t.Error("Metrics not set")
	}
}

func TestNewStateMachine_WithTracer(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create a tracer
	tracer := noop.NewTracerProvider().Tracer("test")

	// Create the state machine with tracer
	sm := NewStateMachine(definition, registry, logger, WithTracer(tracer))

	// Verify that the state machine was created
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}

	// Verify that tracer was set
	if sm.tracer == nil {
		t.Error("Tracer not set")
	}
}

func TestNewStateMachine_WithBothOptions(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create a prometheus registry
	promRegistry := prometheus.NewRegistry()

	// Create a tracer
	tracer := noop.NewTracerProvider().Tracer("test")

	// Create the state machine with both options
	sm := NewStateMachine(definition, registry, logger, WithMetrics(promRegistry), WithTracer(tracer))

	// Verify that the state machine was created
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}

	// Verify that metrics were set
	if sm.metrics == nil {
		t.Error("Metrics not set")
	}

	// Verify that tracer was set
	if sm.tracer == nil {
		t.Error("Tracer not set")
	}
}

func TestNewStateMachine_InvalidDefinition(t *testing.T) {
	// Create an invalid workflow definition (empty states)
	definition := &WorkflowDefinition{
		States: map[string]State{},
	}

	// Create a registry
	registry := NewRegistry()

	// Create a logger
	logger := slog.Default()

	// Create the state machine - this should still work but log an error
	sm := NewStateMachine(definition, registry, logger)

	// With an invalid definition, the state machine should still be created
	// but validation will happen during Trigger
	// Note: NewStateMachine returns nil when the definition is invalid
	// This is the expected behavior
	if sm != nil {
		t.Error("NewStateMachine should return nil for invalid definition")
	}
}

func TestNewStateMachine_NilLogger(t *testing.T) {
	// Create a workflow definition
	definition := &WorkflowDefinition{
		States: map[string]State{
			"start": {
				Name: "start",
			},
		},
	}

	// Create a registry
	registry := NewRegistry()

	// Create the state machine with nil logger
	sm := NewStateMachine(definition, registry, nil)

	// Verify that the state machine was created
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}

	// Verify that a default logger was set
	if sm.logger == nil {
		t.Error("Default logger not set")
	}
}
