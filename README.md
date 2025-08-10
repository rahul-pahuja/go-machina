# GoMachina

[![Go Report Card](https://goreportcard.com/badge/github.com/rahulpahuja/go-machina)](https://goreportcard.com/report/github.com/rahulpahuja/go-machina)
[![GoDoc](https://pkg.go.dev/badge/github.com/rahulpahuja/go-machina)](https://pkg.go.dev/github.com/rahulpahuja/go-machina)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoMachina is a declarative, observable, and extensible state machine engine for building robust, complex workflows in Go. It empowers developers to separate workflow structure from business logic, leading to cleaner, more maintainable, and highly observable applications.

## Why GoMachina?

Managing state in complex processes can be challenging. GoMachina addresses this by providing a robust framework built on four core principles:

-   **Declarative-First Approach**: Define your workflow's entire structure in a clear and simple YAML file. This separation of concerns means your Go code can focus purely on business logic, while the YAML defines the flow, making the workflow easy to visualize and modify without touching the code.

-   **First-Class Observability**: Stop guessing what your system is doing. GoMachina comes with production-grade, built-in observability hooks for structured logging (`slog`), metrics (Prometheus), and tracing (OpenTelemetry). This allows you to monitor, debug, and audit every state transition and action with precision.

-   **Built for Concurrency & Safety**: With a stateless engine that is safe for concurrent use and a guaranteed execution order (`OnLeave` → `Transition Actions` → `OnEnter`), you can build predictable and reliable systems that behave correctly under load.

-   **Powerful & Extensible by Design**: GoMachina is domain-agnostic and designed to be extended. You provide the business logic as simple Go functions, and the engine orchestrates them. It includes built-in support for advanced patterns like "Side Quests" (temporary workflow diversions) and dynamic transitions, allowing you to model even the most complex user journeys.

## Table of Contents

-   [Installation](#installation)
-   [The GoMachina Workflow](#the-gomachina-workflow)
-   [YAML Configuration In-Depth](#yaml-configuration-in-depth)
-   [Implementing Business Logic](#implementing-business-logic)
-   [Putting It All Together](#putting-it-all-together)
-   [Advanced Pattern: Side Quests](#advanced-pattern-side-quests)
-   [API Design & Philosophy](#api-design--philosophy)
-   [Observability](#observability)
-   [For Contributors](#for-contributors)
-   [Roadmap](#roadmap)
-   [License](#license)

## Installation

```bash
go get github.com/rahulpahuja/go-machina
```

## The GoMachina Workflow

Using GoMachina follows a simple, logical lifecycle:

1.  **Define**: Model your entire process as a state machine in a `workflow.yaml` file.
2.  **Implement**: Write standard Go functions for your business logic (Actions) and decision points (Conditions).
3.  **Register**: Map the names from your YAML file to your Go functions in the GoMachina `Registry`.
4.  **Instantiate**: Create an instance of the `StateMachine`, providing it with the loaded definition and registry.
5.  **Execute**: Trigger events to drive the workflow, passing in application data and a context.

## YAML Configuration In-Depth

The `workflow.yaml` file is the heart of GoMachina. It provides a complete, declarative definition of your state machine. Here is an example showcasing many of its features using abstract `A, B, C` states.

```yaml
# The state where the workflow begins.
initialState: A

states:
  # A unique name for the state. This key must match the `name` field below.
  A:
    name: A
    # `onEnter` actions are executed every time this state is entered.
    onEnter:
      - "logEnteringA"
    # `onLeave` actions are executed every time this state is exited.
    onLeave:
      - "logLeavingA"
    transitions:
      # A transition is a link from this state to another, triggered by an `event`.
      - event: "event_to_B"
        # `target` is the state to transition to if conditions pass.
        target: "B"
        # `conditions` are checks that must ALL pass for the transition to occur.
        conditions:
          - "isConditionForB_true"
        # `actions` are executed only during this specific transition.
        actions:
          - "performActionForB"

  B:
    name: B
    onEnter:
      - "logEnteringB"
    transitions:
      - event: "event_to_C"
        target: "C"
        # `autoEvent` immediately triggers the next event in the chain,
        # creating an automated workflow without external triggers.
        autoEvent: "event_to_D"

  C:
    name: C
    onEnter:
      - "logEnteringC"
    transitions:
      - event: "event_to_D"
        target: "D"

  D:
    name: D
    onEnter:
      - "logEnteringD"
```

## Implementing Business Logic

Your Go code provides the implementation for the names defined in the YAML.

```go
package core

import (
    "context"
    "fmt"
    "log/slog"
)

// An ActionFunc executes logic. It can modify and return the data map.
func PerformActionForB(ctx context.Context, data map[string]any) (map[string]any, error) {
    slog.Info("Performing action for state B", "data", data)
    data["action_B_executed_at"] = time.Now()
    return data, nil
}

// A ConditionFunc acts as a guard. It must return true for the transition to proceed.
func IsConditionForB_True(ctx context.Context, data map[string]any) (bool, error) {
    someValue, ok := data["some_key"].(int)
    if !ok {
        return false, nil // Or return an error
    }
    isValid := someValue > 10
    slog.Info("Checked condition for B", "isValid", isValid)
    return isValid, nil
}
```

## Putting It All Together

Here is how you load the definition, register your functions, and run the state machine.

```go
package main

import (
    "context"
    "log"
    "os"
    "log/slog"
    "github.com/rahulpahuja/go-machina/machina"
    "github.com/rahulpahuja/go-machina/core" // Your package with implementations
)

func main() {
    definition, err := machina.LoadWorkflowDefinition("workflow.yaml")
    if err != nil { log.Fatal(err) }

    registry := machina.NewRegistry()
    registry.RegisterCondition("isConditionForB_true", core.IsConditionForB_True)
    registry.RegisterAction("performActionForB", core.PerformActionForB)
    // ... register all other log actions

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    fsm := machina.NewStateMachine(definition, registry, logger)

    ctx := context.Background()
    currentState := definition.InitialState
    currentData := map[string]any{"some_key": 42}
    
    event := "event_to_B"

    // Loop to handle auto-events until a stable state is reached
    for {
        slog.Info("Triggering event", "state", currentState, "event", event)
        result, err := fsm.Trigger(ctx, currentState, event, currentData)
        if err != nil { log.Fatalf("Workflow transition failed: %v", err) }

        currentState = result.NewState
        currentData = result.PersistenceData

        if result.AutoEvent != "" {
            event = result.AutoEvent
            continue
        }
        break
    }
    log.Printf("Workflow finished. Final state: %s", currentState)
}
```

## Advanced Pattern: Side Quests

A "Side Quest" is a temporary diversion from a primary workflow. This powerful pattern allows you to model complex user journeys, such as filling out a sub-form before returning to the main flow.

This is achieved with two core mechanisms:

1.  **The Workflow Stack**: A list of state names, acting as a "breadcrumb trail." This is managed by your actions and stored in the data map, typically under the key `workflow_stack`.
2.  **Dynamic Transition Target**: An action can dynamically set the next state by returning a special `__next_state_override` key in its results. The built-in `__RETURN_TO_PREVIOUS_STATE__` action does this by popping a state from the `workflow_stack`.

Below is a complete example demonstrating this pattern.

**1. The YAML Configuration**

```yaml
# workflow.yaml
states:
  # ... other states like A, B, C ...

  D:
    name: D
    onEnter: ["logAction"]
    transitions:
      - event: "next"
        target: "E"
      # This transition initiates the side quest from state D.
      - event: "start_side_quest"
        target: "B_sharp" # The side quest state.
        actions:
          # This user-defined action pushes the current state ("D") onto the stack.
          - "pushCurrentStateToStack"

  E:
    name: E
    onEnter: ["logAction"]

  # --- Side Quest State --- 
  B_sharp:
    name: B_sharp
    isSideQuest: true # Optional flag to identify this as a side quest.
    onEnter: ["logAction"]
    transitions:
      # This transition completes the side quest and returns to the main flow.
      - event: "complete_and_return"
        # The target is empty because the action will determine it dynamically.
        target: ""
        actions:
          # This special, built-in action pops a state from the workflow stack
          # and sets it as the dynamic target for this transition.
          - "__RETURN_TO_PREVIOUS_STATE__"
```

**2. The Go Implementation**

You only need to implement the action that pushes the state onto the stack. `__RETURN_TO_PREVIOUS_STATE__` is provided by the engine.

```go
// pushCurrentStateToStack adds the current state name to a list in the data map.
func pushCurrentStateToStack(ctx context.Context, data map[string]any) (map[string]any, error) {
    // The FSM passes the source state name in the context.
    sourceState, ok := machina.GetSourceState(ctx)
    if !ok {
        return nil, fmt.Errorf("could not determine source state")
    }

    // Get or create the stack from the data map using the 'workflow_stack' key.
    var stack []string
    if val, ok := data["workflow_stack"]; ok {
        stack = val.([]string)
    }

    // Push the source state onto the stack.
    stack = append(stack, sourceState)
    data["workflow_stack"] = stack
    
    slog.Info("Pushed state to stack", "state", sourceState, "newStack", stack)
    return data, nil
}
```

## API Design & Philosophy

GoMachina is built on a set of core design principles:

-   **Context Propagation**: All operations accept a `context.Context` for timeouts, cancellation, and passing request-scoped data.
-   **Idiomatic Errors**: Errors are handled cleanly, returning `error` values and using `errors.Is` for inspection.
-   **Concurrency Safety**: The FSM engine is stateless and safe for concurrent use. The registry is protected by mutexes.
-   **Extensibility (Strategy Pattern)**: The use of `ActionFunc` and `ConditionFunc` with a central registry allows infinite extension without modifying the core library.
-   **Observability (Observer Pattern)**: The observability hooks for logging, metrics, and tracing allow external systems to monitor the FSM without tight coupling.

## Observability

Inject your own logger, Prometheus registry, and OpenTelemetry tracer upon initialization.

```go
fsm := machina.NewStateMachine(definition, registry, logger,
    machina.WithMetrics(promRegistry),
    machina.WithTracer(tracer),
)
```

-   **Logging**: Provides structured logs with details like `state`, `event`, `transition`, and `duration_ms` for every step.
-   **Metrics**:
    -   `fsm_transitions_total`: Total count of state transitions (labeled by state, event, and target).
    -   `fsm_transition_duration_seconds`: Histogram of transition durations.
    -   `fsm_transition_errors_total`: Total count of errors during transitions.
-   **Tracing**: Creates spans for each transition, allowing you to visualize the workflow in distributed tracing systems.

## For Contributors

We welcome contributions! Here is what you need to know.

### Project Structure
```
/go-machina/
├── /examples/             # Example applications demonstrating features
├── /machina/              # Core library source code
│   ├── definition.go      # Structs for YAML parsing
│   ├── fsm.go             # The state machine engine
│   ├── registry.go        # The action/condition registry
│   ├── metrics.go         # Prometheus metrics implementation
│   ├── loader.go          # YAML/JSON configuration loader
│   ├── validation.go      # Workflow validation logic
│   └── interfaces.go      # Core type definitions (ActionFunc, etc.)
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── README.md              # This file
```

### Development Workflow

This project uses `make` for common development tasks.

-   `make check`: Run all linters and tests.
-   `make test`: Run all tests.
-   `make test-coverage`: Run tests and view coverage.
-   `make lint`: Run the linter.
-   `make fmt`: Format the code.

## Roadmap

-   **State Persistence**: Built-in support for persisting workflow state to databases.
-   **Workflow Visualization**: Tools to generate diagrams from YAML definitions.
-   **Hierarchical State Machines**: Support for nested state machines.
-   **Time-based Transitions**: Trigger transitions after a certain amount of time has passed.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
