package machina

import (
	"os"
	"testing"
)

func TestLoadWorkflowDefinition(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlContent := `
states:
  start:
    name: start
    transitions:
      - event: "validate"
        target: "processOrder"
        conditions:
          - "isUserValid"

  processOrder:
    name: processOrder
    transitions:
      - event: "process"
        target: "complete"
        actions:
          - "chargePayment"
        conditions:
          - "isPaymentSuccess"
      - event: "fail"
        target: "failed"
        actions:
          - "handleFailure"

  complete:
    name: complete
    onEnter:
      - "sendReceipt"

  failed:
    name: failed
`

	tmpfile, err := os.CreateTemp("", "workflow*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load workflow definition
	definition, err := LoadWorkflowDefinition(tmpfile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the loaded definition
	if len(definition.States) != 4 {
		t.Errorf("Expected 4 states, got %d", len(definition.States))
	}

	if definition.States["start"].Name != "start" {
		t.Errorf("Expected first state to be 'start', got %s", definition.States["start"].Name)
	}

	if len(definition.States["start"].Transitions) != 1 {
		t.Errorf("Expected 1 transition in start state, got %d", len(definition.States["start"].Transitions))
	}

	if definition.States["start"].Transitions[0].Event != "validate" {
		t.Errorf("Expected event to be 'validate', got %s", definition.States["start"].Transitions[0].Event)
	}

	if definition.States["start"].Transitions[0].Target != "processOrder" {
		t.Errorf("Expected target state to be 'processOrder', got %s", definition.States["start"].Transitions[0].Target)
	}

	if len(definition.States["start"].Transitions[0].Conditions) != 1 {
		t.Errorf("Expected 1 condition in transition, got %d", len(definition.States["start"].Transitions[0].Conditions))
	}

	if definition.States["start"].Transitions[0].Conditions[0] != "isUserValid" {
		t.Errorf("Expected condition to be 'isUserValid', got %s", definition.States["start"].Transitions[0].Conditions[0])
	}
}

func TestLoadWorkflowDefinition_FileNotFound(t *testing.T) {
	// Try to load a non-existent file
	_, err := LoadWorkflowDefinition("non-existent-file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}
}

func TestLoadWorkflowDefinition_InvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	invalidYAML := `
states:
  start:
    name: start
    transitions:
      event: "validate"
      target: "processOrder"
`

	tmpfile, err := os.CreateTemp("", "invalid-workflow*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(invalidYAML)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Try to load the invalid YAML file
	_, err = LoadWorkflowDefinition(tmpfile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid YAML, got nil")
	}
}
