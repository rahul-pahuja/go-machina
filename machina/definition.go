package machina

// State represents a state in the state machine configuration
type State struct {
	Name        string       `yaml:"name" json:"name"`
	OnEnter     []string     `yaml:"onEnter,omitempty" json:"onEnter,omitempty"`
	OnLeave     []string     `yaml:"onLeave,omitempty" json:"onLeave,omitempty"`
	Transitions []Transition `yaml:"transitions,omitempty" json:"transitions,omitempty"`
}

// Transition represents a transition definition in the configuration
type Transition struct {
	Event      string   `yaml:"event" json:"event"`
	Target     string   `yaml:"target" json:"target"`
	Conditions []string `yaml:"conditions,omitempty" json:"conditions,omitempty"`
	Actions    []string `yaml:"actions,omitempty" json:"actions,omitempty"`
	AutoEvent  string   `yaml:"autoEvent,omitempty" json:"autoEvent,omitempty"` // Event to automatically fire after transition
}

// WorkflowDefinition represents the entire workflow configuration
type WorkflowDefinition struct {
	States map[string]State `yaml:"states" json:"states"`
}
