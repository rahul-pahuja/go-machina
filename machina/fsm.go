//go:generate go run ../tools/generate_mocks.go

package machina

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TransitionResult holds all the successful outcomes of a Trigger event.
type TransitionResult struct {
	NewState        string
	AutoEvent       string
	PersistenceData map[string]any
}

// StateMachine represents the finite state machine
type StateMachine struct {
	definition *WorkflowDefinition
	registry   *Registry
	logger     *slog.Logger
	metrics    *Metrics
	tracer     trace.Tracer
}

// StateMachineOption is a function that configures a StateMachine
type StateMachineOption func(*StateMachine)

// WithMetrics configures the StateMachine with Prometheus metrics
func WithMetrics(reg prometheus.Registerer) StateMachineOption {
	return func(sm *StateMachine) {
		sm.metrics = NewMetrics(reg)
	}
}

// WithTracer configures the StateMachine with OpenTelemetry tracing
func WithTracer(tracer trace.Tracer) StateMachineOption {
	return func(sm *StateMachine) {
		sm.tracer = tracer
	}
}

// NewStateMachine creates a new state machine instance
func NewStateMachine(definition *WorkflowDefinition, registry *Registry, logger *slog.Logger, opts ...StateMachineOption) *StateMachine {
	if logger == nil {
		logger = slog.Default()
	}

	// Validate the workflow definition
	if err := definition.Validate(); err != nil {
		logger.Error("Invalid workflow definition", "error", err)
		return nil
	}

	// Register the predefined RETURN_TO_PREVIOUS_STATE action
	registry.RegisterAction("__RETURN_TO_PREVIOUS_STATE__", ReturnToPreviousStateAction)

	sm := &StateMachine{
		definition: definition,
		registry:   registry,
		logger:     logger,
		tracer:     otel.Tracer("gomachina"),
		// Initialize with no-op metrics by default
		metrics: NewMetrics(nil),
	}

	// Apply options
	for _, opt := range opts {
		opt(sm)
	}

	return sm
}

// Trigger processes a single event and causes a state transition
func (sm *StateMachine) Trigger(ctx context.Context, currentState string, event string, payload map[string]any) (*TransitionResult, error) {
	startTime := time.Now()

	// Create a span for tracing
	ctx, span := sm.tracer.Start(ctx, "fsm.transition",
		trace.WithAttributes(
			attribute.String("fsm.current_state", currentState),
			attribute.String("fsm.event", event),
		))
	defer span.End()

	// Find the current state definition
	stateDef, err := sm.getStateDefinition(currentState)
	if err != nil {
		err = fmt.Errorf("failed to get state definition for %s: %w", currentState, err)
		sm.recordTransitionError(currentState, event, "state_not_found", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	sm.logger.Info("Processing event", "state", currentState, "event", event, "payload", payload)

	// Find the transition for the event
	transition, err := sm.getTransitionForEvent(stateDef, event, ctx, payload)
	if err != nil {
		err = fmt.Errorf("no valid transition found for event %s in state %s: %w", event, currentState, err)
		sm.recordTransitionError(currentState, event, "transition_not_found", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String("fsm.target_state", transition.Target),
		attribute.StringSlice("fsm.conditions", transition.Conditions),
		attribute.StringSlice("fsm.actions", transition.Actions),
	)

	sm.logger.Info("Found transition", "event", event, "target", transition.Target, "conditions", transition.Conditions, "actions", transition.Actions)

	// Initialize persistenceData as a copy of the payload to avoid modifying the original
	persistenceData := make(map[string]any)
	for k, v := range payload {
		persistenceData[k] = v
	}

	// Check all conditions for the transition
	if err := sm.executeConditions(ctx, currentState, event, transition, payload); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Execute transition actions (proposed new order)
	if err := sm.executeTransitionActions(ctx, currentState, event, transition.Actions, payload, persistenceData); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Check for dynamic transition target override
	nextStateOverride, hasOverride := persistenceData["__next_state_override"]
	if hasOverride {
		if overrideStr, ok := nextStateOverride.(string); ok && overrideStr != "" {
			transition.Target = overrideStr
			span.SetAttributes(attribute.String("fsm.dynamic_target", overrideStr))
			sm.logger.Info("Dynamic transition target override", "from", transition.Target, "to", overrideStr)
			// Clear the override value so it doesn't affect future transitions
			delete(persistenceData, "__next_state_override")
		}
	}

	// Execute OnLeave actions for the current state
	if err := sm.executeOnLeaveActions(ctx, currentState, event, stateDef.OnLeave, payload, persistenceData); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Execute OnEnter actions for the target state
	targetStateDef, err := sm.getStateDefinition(transition.Target)
	if err != nil {
		err = fmt.Errorf("failed to get target state definition for %s: %w", transition.Target, err)
		sm.recordTransitionError(currentState, event, "target_state_not_found", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := sm.executeOnEnterActions(ctx, currentState, event, transition.Target, targetStateDef.OnEnter, payload, persistenceData); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Record successful transition metrics
	duration := time.Since(startTime).Seconds()
	if sm.metrics != nil {
		sm.metrics.TransitionsTotal.WithLabelValues(currentState, transition.Target, event).Inc()
		sm.metrics.TransitionDuration.WithLabelValues(currentState, transition.Target, event).Observe(duration)

		// Record auto transition if applicable
		if transition.AutoEvent != "" {
			sm.metrics.AutoTransitionsTotal.WithLabelValues(currentState, transition.Target, event).Inc()
		}
	}

	sm.logger.Info("Transition completed", "from", currentState, "to", transition.Target, "event", event, "duration_seconds", duration)
	span.SetAttributes(
		attribute.String("fsm.new_state", transition.Target),
		attribute.Float64("fsm.duration_seconds", duration),
	)

	return &TransitionResult{
		NewState:        transition.Target,
		AutoEvent:       transition.AutoEvent,
		PersistenceData: persistenceData,
	}, nil
}

// GetAutoEventForTransition returns the auto event for a transition, if any
func (sm *StateMachine) GetAutoEventForTransition(fromState, event string) (string, error) {
	stateDef, err := sm.getStateDefinition(fromState)
	if err != nil {
		return "", fmt.Errorf("failed to get state definition for %s: %w", fromState, err)
	}

	// Use a background context and empty payload for auto event lookup
	transition, err := sm.getTransitionForEvent(stateDef, event, context.Background(), map[string]any{})
	if err != nil {
		return "", fmt.Errorf("no valid transition found for event %s in state %s: %w", event, fromState, err)
	}

	return transition.AutoEvent, nil
}

// getStateDefinition finds a state definition by name
func (sm *StateMachine) getStateDefinition(name string) (*State, error) {
	state, exists := sm.definition.States[name]
	if !exists {
		return nil, fmt.Errorf("state %s not found", name)
	}
	return &state, nil
}

// getTransitionForEvent finds the transition for a specific event in a state
// For conditional transitions, it evaluates conditions and returns the first matching transition
func (sm *StateMachine) getTransitionForEvent(state *State, event string, ctx context.Context, payload map[string]any) (*Transition, error) {
	var matchingTransitions []Transition
	
	// Collect all transitions for the event
	for _, transition := range state.Transitions {
		if transition.Event == event {
			matchingTransitions = append(matchingTransitions, transition)
		}
	}
	
	if len(matchingTransitions) == 0 {
		return nil, fmt.Errorf("no transition found for event %s", event)
	}
	
	// If only one transition, return it directly
	if len(matchingTransitions) == 1 {
		return &matchingTransitions[0], nil
	}
	
	// Multiple transitions - evaluate conditions to find the first matching one
	for _, transition := range matchingTransitions {
		// If no conditions, this is a match
		if len(transition.Conditions) == 0 {
			return &transition, nil
		}
		
		// Evaluate all conditions
		allConditionsMet := true
		for _, conditionName := range transition.Conditions {
			condition, err := sm.registry.GetCondition(conditionName)
			if err != nil {
				return nil, fmt.Errorf("failed to get condition %s: %w", conditionName, err)
			}
			
			ok, err := condition(ctx, payload)
			if err != nil {
				return nil, fmt.Errorf("condition %s failed: %w", conditionName, err)
			}
			
			if !ok {
				allConditionsMet = false
				break
			}
		}
		
		// If all conditions are met, this is our transition
		if allConditionsMet {
			return &transition, nil
		}
	}
	
	return nil, fmt.Errorf("no transition found for event %s with matching conditions", event)
}

// mergeData merges two data maps
func (sm *StateMachine) mergeData(original, updates map[string]any) map[string]any {
	// Merge the maps
	result := make(map[string]any)
	for k, v := range original {
		result[k] = v
	}
	for k, v := range updates {
		result[k] = v
	}

	return result
}

// executeConditions checks all conditions for a transition
func (sm *StateMachine) executeConditions(ctx context.Context, currentState, event string, transition *Transition, payload map[string]any) error {
	for _, conditionName := range transition.Conditions {
		condition, err := sm.registry.GetCondition(conditionName)
		if err != nil {
			err = fmt.Errorf("failed to get condition %s: %w", conditionName, err)
			sm.recordTransitionError(currentState, event, "condition_not_found", err)
			return err
		}

		sm.logger.Info("Evaluating condition", "condition", conditionName)
		ok, err := condition(ctx, payload)
		if err != nil {
			err = fmt.Errorf("condition %s failed: %w", conditionName, err)
			sm.recordTransitionError(currentState, event, "condition_error", err)
			return err
		}

		if !ok {
			err = fmt.Errorf("condition %s evaluated to false", conditionName)
			sm.recordTransitionError(currentState, event, "condition_failed", err)
			sm.logger.Info("Condition evaluated to false", "condition", conditionName)
			return err
		}

		sm.logger.Info("Condition passed", "condition", conditionName)
	}
	return nil
}

// executeTransitionActions executes transition actions
func (sm *StateMachine) executeTransitionActions(ctx context.Context, currentState, event string, actions []string, payload map[string]any, persistenceData map[string]any) error {
	for _, actionName := range actions {
		action, err := sm.registry.GetAction(actionName)
		if err != nil {
			err = fmt.Errorf("failed to get transition action %s: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "transition_action_not_found", err)
			return err
		}

		sm.logger.Info("Executing transition action", "action", actionName)
		result, err := action(ctx, payload)
		if err != nil {
			err = fmt.Errorf("transition action %s failed: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "transition_action_error", err)
			return err
		}

		// Update persistenceData with result
		if result != nil {
			for k, v := range result {
				persistenceData[k] = v
			}
			sm.logger.Info("Transition action updated persistenceData", "action", actionName, "updates", result)
		}
	}
	return nil
}

// executeOnLeaveActions executes OnLeave actions for the current state
func (sm *StateMachine) executeOnLeaveActions(ctx context.Context, currentState, event string, actions []string, payload map[string]any, persistenceData map[string]any) error {
	for _, actionName := range actions {
		action, err := sm.registry.GetAction(actionName)
		if err != nil {
			err = fmt.Errorf("failed to get OnLeave action %s: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "onleave_action_not_found", err)
			return err
		}

		sm.logger.Info("Executing OnLeave action", "action", actionName)
		result, err := action(ctx, payload)
		if err != nil {
			err = fmt.Errorf("OnLeave action %s failed: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "onleave_action_error", err)
			return err
		}

		// Update persistenceData with result
		if result != nil {
			for k, v := range result {
				persistenceData[k] = v
			}
			sm.logger.Info("OnLeave action updated persistenceData", "action", actionName, "updates", result)
		}
	}
	return nil
}

// executeOnEnterActions executes OnEnter actions for the target state
func (sm *StateMachine) executeOnEnterActions(ctx context.Context, currentState, event, targetState string, actions []string, payload map[string]any, persistenceData map[string]any) error {
	for _, actionName := range actions {
		action, err := sm.registry.GetAction(actionName)
		if err != nil {
			err = fmt.Errorf("failed to get OnEnter action %s: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "onenter_action_not_found", err)
			return err
		}

		sm.logger.Info("Executing OnEnter action", "action", actionName)
		result, err := action(ctx, payload)
		if err != nil {
			err = fmt.Errorf("OnEnter action %s failed: %w", actionName, err)
			sm.recordTransitionError(currentState, event, "onenter_action_error", err)
			return err
		}

		// Update persistenceData with result
		if result != nil {
			for k, v := range result {
				persistenceData[k] = v
			}
			sm.logger.Info("OnEnter action updated persistenceData", "action", actionName, "updates", result)
		}
	}
	return nil
}

// recordTransitionError records a transition error in metrics
func (sm *StateMachine) recordTransitionError(fromState, event, errorType string, err error) {
	if sm.metrics != nil {
		sm.metrics.TransitionErrors.WithLabelValues(fromState, event, errorType).Inc()
	}
}

// ReturnToPreviousStateAction is a predefined action that pops the top state from the WorkflowStack
// and returns it as the __next_state_override
func ReturnToPreviousStateAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	// Get the workflow stack from the context
	workflowStack, ok := data["WorkflowStack"].([]string)
	if !ok || len(workflowStack) == 0 {
		return nil, fmt.Errorf("workflow stack not found or empty")
	}

	// Pop the top state from the stack
	returnState := workflowStack[len(workflowStack)-1]
	workflowStack = workflowStack[:len(workflowStack)-1]

	// Return the popped state as the next state override and updated stack
	return map[string]any{
		"__next_state_override": returnState,
		"WorkflowStack":         workflowStack,
	}, nil
}
