package machina

import (
	"fmt"
	"sync"
)

// Registry holds mappings of condition and action implementations
type Registry struct {
	conditions map[string]ConditionFunc
	actions    map[string]ActionFunc
	mu         sync.RWMutex
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		conditions: make(map[string]ConditionFunc),
		actions:    make(map[string]ActionFunc),
	}
}

// RegisterCondition registers a condition function
func (r *Registry) RegisterCondition(name string, condition ConditionFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.conditions[name]; exists {
		return fmt.Errorf("condition %s already registered", name)
	}

	r.conditions[name] = condition
	return nil
}

// RegisterAction registers an action function
func (r *Registry) RegisterAction(name string, action ActionFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.actions[name]; exists {
		return fmt.Errorf("action %s already registered", name)
	}

	r.actions[name] = action
	return nil
}

// GetCondition retrieves a condition function by name
func (r *Registry) GetCondition(name string) (ConditionFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if condition, exists := r.conditions[name]; exists {
		return condition, nil
	}

	return nil, fmt.Errorf("condition %s not found", name)
}

// GetAction retrieves an action function by name
func (r *Registry) GetAction(name string) (ActionFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if action, exists := r.actions[name]; exists {
		return action, nil
	}

	return nil, fmt.Errorf("action %s not found", name)
}
