package machina

import (
	"fmt"
)

// Validate checks if the workflow definition is valid
func (wd *WorkflowDefinition) Validate() error {
	if len(wd.States) == 0 {
		return fmt.Errorf("workflow must have at least one state")
	}

	// Validate initial state if specified
	if wd.InitialState != "" {
		if _, exists := wd.States[wd.InitialState]; !exists {
			return fmt.Errorf("initialState %s not found in states", wd.InitialState)
		}
	}

	// Validate each state
	for name, state := range wd.States {
		if name != state.Name {
			return fmt.Errorf("state key %s does not match state name %s", name, state.Name)
		}

		if err := state.Validate(); err != nil {
			return fmt.Errorf("invalid state %s: %w", state.Name, err)
		}
	}

	return nil
}

// Validate checks if the state is valid
func (s *State) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("state must have a name")
	}

	// Validate transitions
	for _, transition := range s.Transitions {
		if err := transition.Validate(); err != nil {
			return fmt.Errorf("invalid transition for event %s: %w", transition.Event, err)
		}
	}

	return nil
}

// Validate checks if the transition is valid
func (t *Transition) Validate() error {
	if t.Event == "" {
		return fmt.Errorf("transition must have an event")
	}

	if t.Target == "" {
		return fmt.Errorf("transition must have a target state")
	}

	return nil
}
