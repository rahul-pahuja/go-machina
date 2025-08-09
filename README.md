# GoMachina

A declarative state machine engine for Go applications.

[![Go Report Card](https://goreportcard.com/badge/github.com/rahulpahuja/go-machina)](https://goreportcard.com/report/github.com/rahulpahuja/go-machina)
[![GoDoc](https://pkg.go.dev/badge/github.com/rahulpahuja/go-machina)](https://pkg.go.dev/github.com/rahulpahuja/go-machina)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

GoMachina is a powerful, declarative, and observable state machine library for orchestrating complex business workflows with safety, clarity, and extensibility. It allows you to define your business logic as a finite state machine using YAML/JSON configuration files, while implementing the actual business logic in Go code.

## Features

- **Declarative-First**: Define state machine structure in YAML/JSON configuration files
- **High Observability**: Built-in structured logging, metrics (Prometheus), and tracing (OpenTelemetry)
- **Safety & Predictability**: Safe execution order (`OnLeave` → `Transition Actions` → `OnEnter`) and concurrent use
- **Extensible & Decoupled**: Generic, domain-agnostic design with user-provided implementations
- **Auto-Events**: Support for automatic event triggering for seamless workflow execution
- **Validation**: Built-in validation for workflow definitions
- **Context Support**: Full context propagation for cancellation and timeouts

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Configuration](#configuration)
- [Implementation](#implementation)
- [Usage](#usage)
- [Observability](#observability)
- [Advanced Features](#advanced-features)
- [Development](#development)
- [Testing](#testing)
- [License](#license)

## Installation

```bash
go get github.com/rahulpahuja/go-machina
```

## Quick Start

1. Create a workflow definition (`workflow.yaml`):

```yaml
initialState: start
states:
  start:
    name: start
    onEnter:
      - "logStart"
    transitions:
      - event: "validate"
        target: "processOrder"
        conditions:
          - "isUserValid"

  processOrder:
    name: processOrder
    onEnter:
      - "logProcessing"
    transitions:
      - event: "process"
        target: "complete"
        actions:
          - "chargePayment"
        conditions:
          - "isPaymentSuccess"

  complete:
    name: complete
    onEnter:
      - "logCompletion"
      - "sendReceipt"
```

2. Implement your conditions and actions:

```go
// Condition function type
type ConditionFunc func(ctx context.Context, data map[string]any) (bool, error)

// Action function type
type ActionFunc func(ctx context.Context, data map[string]any) (map[string]any, error)

// Example condition implementation
func IsUserValidCondition(ctx context.Context, data map[string]any) (bool, error) {
    user, ok := data["user"].(map[string]any)
    if !ok {
        return false, fmt.Errorf("user data not found")
    }
    
    userID, ok := user["id"].(string)
    if !ok || userID == "" {
        return false, nil
    }
    
    // In a real implementation, you would check against a database or service
    return true, nil
}

// Example action implementation
func ChargePaymentAction(ctx context.Context, data map[string]any) (map[string]any, error) {
    amount, ok := data["amount"].(float64)
    if !ok {
        return nil, fmt.Errorf("amount not found or invalid")
    }
    
    // In a real implementation, you would integrate with a payment provider
    log.Printf("Charging payment of $%.2f", amount)
    
    return map[string]any{
        "paymentStatus": "charged",
        "chargedAt":     time.Now(),
    }, nil
}
```

3. Register your implementations and run the state machine:

```go
// Load workflow definition
definition, err := machina.LoadWorkflowDefinition("workflow.yaml")
if err != nil {
    log.Fatal(err)
}

// Create registry and register conditions/actions
registry := machina.NewRegistry()
registry.RegisterCondition("isUserValid", IsUserValidCondition)
registry.RegisterAction("chargePayment", ChargePaymentAction)
// ... register other conditions and actions

// Create state machine with observability options
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
fsm := machina.NewStateMachine(definition, registry, logger,
    machina.WithMetrics(prometheus.DefaultRegisterer),
    machina.WithTracer(otel.Tracer("my-app")))

// Execute workflow
ctx := context.Background()
data := map[string]any{
    "orderId": "order-123",
    "user": map[string]any{
        "id":    "user-456",
        "email": "user@example.com",
    },
    "amount": 99.99,
}

// Handle auto-events in a loop
currentState := definition.InitialState
currentData := data

for {
    // Determine the next event based on business logic
    var event string
    switch currentState {
    case "start":
        event = "validate"
    case "processOrder":
        event = "process"
    default:
        // No more events to process
        break
    }
    
    if event == "" {
        break
    }
    
    // Trigger the event
    result, err := fsm.Trigger(ctx, currentState, event, currentData)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update state and data
    currentState = result.NewState
    currentData = result.PersistenceData
    
    // Handle auto-events
    if result.AutoEvent != "" {
        event = result.AutoEvent
        continue
    }
    
    // Break if we've reached a terminal state
    if currentState == "complete" || currentState == "failed" {
        break
    }
}

fmt.Printf("Workflow completed. Final state: %s\n", currentState)
```

## Core Concepts

### State Machine
A state machine is a mathematical model of computation that represents the behavior of a system. It consists of:
- **States**: Discrete conditions or situations in which the system can exist
- **Transitions**: Movements from one state to another based on events
- **Events**: Triggers that cause transitions between states
- **Actions**: Operations performed during transitions or when entering/leaving states
- **Conditions**: Logical expressions that must be satisfied for transitions to occur

### Workflow Definition
The workflow definition is a declarative representation of your business process in YAML or JSON format. It defines:
- The initial state of the workflow
- All possible states and their properties
- Transitions between states and the events that trigger them
- Conditions that must be met for transitions
- Actions to be executed during transitions or state changes

### Registry
The registry is a central component that maps condition and action names (as defined in the workflow) to their actual implementations in code. This allows for loose coupling between the workflow definition and the business logic.

### Transition Result
When a transition is triggered, the state machine returns a `TransitionResult` struct containing:
- `NewState`: The state the machine transitioned to
- `AutoEvent`: An optional event that should be automatically triggered next
- `PersistenceData`: The data that should be persisted after the transition

## Configuration

### YAML Structure

```yaml
initialState: start  # Optional: The initial state of the workflow
states:
  stateName:
    name: stateName  # Must match the key
    onEnter:         # Optional: Actions to execute when entering this state
      - actionName1
      - actionName2
    onLeave:         # Optional: Actions to execute when leaving this state
      - actionName3
    transitions:
      - event: eventName     # The event that triggers this transition
        target: targetState  # The state to transition to
        conditions:          # Optional: Conditions that must be met
          - conditionName1
          - conditionName2
        actions:             # Optional: Actions to execute during transition
          - actionName4
        autoEvent: autoEventName  # Optional: Event to automatically trigger
```

### Example with All Features

```yaml
initialState: draft
states:
  draft:
    name: draft
    onEnter:
      - "logDraftCreated"
    transitions:
      - event: "submit"
        target: "review"
        conditions:
          - "isValidSubmission"
        actions:
          - "saveSubmission"
  
  review:
    name: review
    onEnter:
      - "logReviewStarted"
      - "notifyReviewers"
    transitions:
      - event: "approve"
        target: "approved"
        actions:
          - "logApproval"
      - event: "reject"
        target: "rejected"
        actions:
          - "logRejection"
          - "notifySubmitter"
        autoEvent: "notify"
  
  approved:
    name: approved
    onEnter:
      - "logApproved"
      - "publishContent"
  
  rejected:
    name: rejected
    onEnter:
      - "logRejected"
```

## Implementation

### Condition Functions

Conditions are functions that determine whether a transition should occur. They must have the signature:

```go
type ConditionFunc func(ctx context.Context, data map[string]any) (bool, error)
```

Example:
```go
func IsPaymentAmountValid(ctx context.Context, data map[string]any) (bool, error) {
    amount, ok := data["amount"].(float64)
    if !ok {
        return false, fmt.Errorf("amount not found")
    }
    
    // Payment must be at least $1
    return amount >= 1.0, nil
}
```

### Action Functions

Actions are functions that perform business logic during transitions. They must have the signature:

```go
type ActionFunc func(ctx context.Context, data map[string]any) (map[string]any, error)
```

Example:
```go
func SendEmailNotification(ctx context.Context, data map[string]any) (map[string]any, error) {
    email, ok := data["email"].(string)
    if !ok {
        return nil, fmt.Errorf("email not found")
    }
    
    subject, ok := data["subject"].(string)
    if !ok {
        return nil, fmt.Errorf("subject not found")
    }
    
    message, ok := data["message"].(string)
    if !ok {
        return nil, fmt.Errorf("message not found")
    }
    
    // In a real implementation, you would send the email
    log.Printf("Sending email to %s: %s", email, subject)
    
    return map[string]any{
        "emailSentAt": time.Now(),
        "emailStatus": "sent",
    }, nil
}
```

### Lifecycle Hooks

GoMachina supports three types of lifecycle hooks:
1. **OnEnter**: Executed when entering a state
2. **Transition Actions**: Executed during a transition
3. **OnLeave**: Executed when leaving a state

The execution order is:
1. Source state's OnLeave actions
2. Transition actions
3. Target state's OnEnter actions

## Usage

### Basic Usage

```bash
# Run the example server
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build the binary
make build

# Install dependencies
make deps

# Clean build artifacts
make clean
```

### Advanced Usage Patterns

#### 1. Context Integration

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Trigger transition with context
result, err := fsm.Trigger(ctx, currentState, event, data)
```

#### 2. Guard Conditions

```go
// Add runtime guard conditions
result, err := fsm.Trigger(ctx, currentState, event, data, 
    func(ctx context.Context, data map[string]any) (bool, error) {
        // Custom guard condition
        return data["allowTransition"] == true, nil
    })
```

#### 3. Error Handling

```go
result, err := fsm.Trigger(ctx, currentState, event, data)
if err != nil {
    // Handle different types of errors
    switch {
    case errors.Is(err, machina.ErrStateNotFound):
        log.Printf("State not found: %v", err)
    case errors.Is(err, machina.ErrTransitionNotFound):
        log.Printf("Transition not found: %v", err)
    case errors.Is(err, machina.ErrConditionFailed):
        log.Printf("Condition failed: %v", err)
    default:
        log.Printf("Unexpected error: %v", err)
    }
    return
}
```

## Observability

### Logging

GoMachina uses structured logging with the `log/slog` package. You can provide your own logger:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
fsm := machina.NewStateMachine(definition, registry, logger)
```

### Metrics

Metrics are exposed via Prometheus:

```go
// Enable metrics
promRegistry := prometheus.NewRegistry()
fsm := machina.NewStateMachine(definition, registry, logger,
    machina.WithMetrics(promRegistry))

// Register with default Prometheus registry if needed
prometheus.MustRegister(promRegistry)
```

Available metrics:
- `fsm_transitions_total`: Count of state transitions
- `fsm_transition_duration_seconds`: Duration of transitions
- `fsm_auto_transitions_total`: Count of auto-transitions
- `fsm_transition_errors_total`: Count of transition errors

### Tracing

Distributed tracing is supported via OpenTelemetry:

```go
// Enable tracing
tracer := otel.Tracer("my-service")
fsm := machina.NewStateMachine(definition, registry, logger,
    machina.WithTracer(tracer))
```

## Advanced Features

### Auto-Events

Auto-events allow for seamless workflow execution without external triggering:

```yaml
transitions:
  - event: "approve"
    target: "approved"
    autoEvent: "notify"  # Automatically trigger "notify" event
```

In code:
```go
for {
    result, err := fsm.Trigger(ctx, currentState, event, data)
    if err != nil {
        // Handle error
        break
    }
    
    currentState = result.NewState
    data = result.PersistenceData
    
    // Auto-trigger next event if specified
    if result.AutoEvent != "" {
        event = result.AutoEvent
        continue
    }
    
    break
}
```

### Validation

Workflow definitions are automatically validated:
- All states must have matching names
- Transitions must have events and targets
- Initial state must exist if specified
- Referenced conditions and actions must be registered

### Concurrency Safety

The state machine is safe for concurrent use:
- Registry uses mutexes to protect shared state
- State machine engine is stateless (state is passed as parameters)
- All public APIs are safe for concurrent access

## Development

### Prerequisites

- Go 1.21 or higher
- Make

### Make Commands

```bash
# Install dependencies
make deps

# Build the project
make build

# Run the example server
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Lint the code
make lint

# Format the code
make fmt

# Clean build artifacts
make clean

# Generate mocks
make generate

# Run all checks (test, lint, etc.)
make check
```

### Project Structure

```
/go-machina/
├── /cmd/
│   └── /server/           # Example server application
├── /configs/              # Configuration files
│   └── workflow.yaml      # Example workflow definition
├── /internal/             # Internal packages
│   └── /core/             # Example implementations
├── /machina/              # Core library
│   ├── definition.go      # Workflow definition types
│   ├── fsm.go             # State machine implementation
│   ├── registry.go        # Condition/action registry
│   ├── metrics.go         # Observability metrics
│   ├── loader.go          # Configuration loader
│   ├── validation.go      # Workflow validation
│   └── interfaces.go      # Type definitions
├── /tools/                # Development tools
├── Makefile               # Build automation
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── README.md              # This file
└── gomachina_blueprint.md # Project blueprint
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run specific package tests
make test-unit
make test-integration
```

### Test Structure

```
/machina/
├── *_test.go              # Unit tests
├── *_benchmark_test.go    # Benchmarks
├── *_error_test.go        # Error condition tests
└── /cmd/server/
    └── integration_test.go # Integration tests
```

### Writing Tests

Example test for a condition:

```go
func TestIsUserValidCondition(t *testing.T) {
    ctx := context.Background()
    
    // Test valid user
    data := map[string]any{
        "user": map[string]any{
            "id": "user123",
        },
    }
    
    valid, err := IsUserValidCondition(ctx, data)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if !valid {
        t.Error("Expected user to be valid")
    }
    
    // Test invalid user
    invalidData := map[string]any{
        "user": map[string]any{
            "id": "",
        },
    }
    
    valid, err = IsUserValidCondition(ctx, invalidData)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if valid {
        t.Error("Expected user to be invalid")
    }
}
```

## Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
