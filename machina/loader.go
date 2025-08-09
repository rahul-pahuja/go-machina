package machina

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadWorkflowDefinition loads a workflow definition from a YAML file
func LoadWorkflowDefinition(filePath string) (*WorkflowDefinition, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var definition WorkflowDefinition
	definition.States = make(map[string]State)

	if err := yaml.Unmarshal(data, &definition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &definition, nil
}
