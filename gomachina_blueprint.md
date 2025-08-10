# Project Blueprint: `GoMachina` - A Declarative State Machine Engine

---

## 1. Project Vision & Goals

### 1.1. Project Name
`GoMachina`

### 1.2. Vision
To provide the Go ecosystem with a powerful, declarative, and observable state machine library for orchestrating complex business workflows with safety, clarity, and extensibility.

### 1.3. Core Principles
**Four unbreakable rules:**
1. **Declarative-First:** State machine structure (states, transitions) is defined in a static configuration file (YAML/JSON). Code implements business logic; configuration defines the flow.
2. **High Observability:** Production-grade observability (structured logging, metrics, tracing) is a built-in feature.
3. **Safety & Predictability:** Safe execution order (`Condition` → `Action`), safe for concurrent use.
4. **Extensible & Decoupled:** Completely generic, domain-agnostic, extended via user-provided implementations.

---

## 2. Core Concepts

| Concept       | Description |
|---------------|-------------|
| **State**     | Represents a discrete step in the workflow. |
| **Transition**| Defines movement from one state to another when a condition passes. |
| **Condition** | Logic that determines if a transition should occur. |
| **Action**    | Business logic executed within a state. |
| **Registry**  | Holds mappings of condition and action implementations. |
| **InitialState** | The starting state of the workflow, defined in the configuration. |
| **TransitionResult** | The result of a state transition, containing the new state, any auto-event, and persistence data. |

---

## 3. High-Level Architecture
```
+---------------------+
|   Config Loader     |
| (YAML/JSON Parser)  |
+---------------------+
           |
           v
+---------------------+       +-------------------+
|  State Machine FSM  | ----> | Observability      |
| (Transitions Engine)|       | (Logs/Metrics)    |
+---------------------+       +-------------------+
     |          |
     v          v
Conditions   Actions
```

---

## 4. Go Package Structure
```
/go-machina/
├── /cmd/
│   └── /server/main.go
├── /configs/
│   └── workflow.yaml
├── /internal/
│   ├── api/          # HTTP handlers and routing
│   ├── core/         # User-provided actions & conditions
│   └── data/         # Data access logic
├── /machina/         # GoMachina core library
│   ├── definition.go
│   ├── fsm.go
│   ├── registry.go
│   ├── metrics.go
│   ├── loader.go
│   ├── validation.go
│   └── interfaces.go
└── go.mod
```

---

## 5. Interfaces
```go
// ConditionFunc defines the function signature for evaluating transition conditions
type ConditionFunc func(ctx context.Context, data map[string]any) (bool, error)

// ActionFunc defines the function signature for executing state actions
// It returns a map of updated data and an error
type ActionFunc func(ctx context.Context, data map[string]any) (map[string]any, error)
```

---

## 6. Example Configuration (YAML)
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
      - event: "fail"
        target: "failed"
        actions:
          - "handleFailure"

  complete:
    name: complete
    onEnter:
      - "logCompletion"
      - "sendReceipt"

  failed:
    name: failed
    onEnter:
      - "logFailure"
```

---

## 7. Error Handling & Recovery
- **Invalid State Transitions:** Logged and halted.
- **Condition Failures:** Gracefully fail transition, log reason.
- **Action Failures:** Support retries with exponential backoff.
- **Observability Hooks:** For debugging and auditing.

---

## 8. Observability
- **Logging:** Structured logs per state transition.
- **Metrics:**
  - Counter: Number of transitions per state.
  - Histogram: Execution time for actions.
- **Tracing:** Optional OpenTelemetry integration.

---

## 9. Sample Use Case — Payment Workflow
1. **Config File:** Defines `initialState: start` and states `start → processOrder → complete`.
2. **Conditions:**
   - `isUserValid`
   - `isPaymentSuccess`
3. **Actions:**
   - `chargePayment`
   - `sendReceipt`
4. **Execution:** FSM processes events sequentially using `fsm.Trigger()`, transitions only on passing conditions. The result is returned as a `TransitionResult` struct containing the new state, any auto-event, and persistence data.

---

## 10. Extensibility Guidelines
- Implement `ConditionFunc` and `ActionFunc` function types for custom logic.
- Register them in the FSM registry using `RegisterCondition` and `RegisterAction`.
- Use `WithMetrics` and `WithTracer` options to customize observability.
- Access the `InitialState` from the workflow definition to start workflows without hardcoding the initial state.
- Handle `AutoEvent` from the `TransitionResult` to implement auto-trigger loops.

---

## 11. Design Patterns

### 11.1. Strategy Pattern
The `Registry` uses the Strategy pattern to allow users to provide their own implementations of `Condition` and `Action` interfaces. This enables the library to be domain-agnostic while still providing powerful functionality.

### 11.2. State Pattern
The core FSM engine implements the State pattern, where the behavior of the system changes based on its current state. Each state encapsulates its own transitions and behavior.

### 11.3. Observer Pattern
The observability features (logging, metrics, tracing) implement the Observer pattern, allowing external systems to monitor the FSM without coupling to its internal implementation.

### 11.4. Decorator Pattern
The optional metrics and tracing capabilities are implemented using the Decorator pattern, wrapping the core functionality with additional observability features without modifying the core implementation.

---

## 12. Implementation Details

### 12.1. State Machine Engine
The core `StateMachine` struct is responsible for:
- Processing events to find a matching transition for the current state.
- Validating the transition by executing all registered `ConditionFunc`s.
- Executing registered `ActionFunc`s for `onLeave` (from the source state) and for the transition itself.
- Collecting and merging persistence data returned from all executed actions.
- **Dynamically determining the next state:** It first checks if the collected persistence data contains a special `__next_state_override` key. If so, it uses that as the target state. Otherwise, it uses the static `target` from the workflow definition.
- Executing `onEnter` actions for the newly determined target state.
- Returning the final result, including the new state and all persistence data, in a `TransitionResult` struct.

### 12.2. Configuration Loader
The configuration loader is responsible for:
- Parsing YAML/JSON workflow definitions
- Validating the structure and consistency of the workflow
- Creating the in-memory representation of the state machine

### 12.3. Registry
The registry is responsible for:
- Storing mappings of condition and action names to implementations
- Providing thread-safe access to these implementations
- Allowing users to register their custom implementations

### 12.4. Internal Actions & Side Quests
To support advanced patterns like "Side Quests," the `GoMachina` engine provides built-in, pre-registered actions that can be used in any workflow definition.

- A key example is `__RETURN_TO_PREVIOUS_STATE__`. This action is designed to be used in a transition to dynamically return to a previous state stored in a `WorkflowStack` (which is part of the application's context data).
- These internal actions are automatically made available in the `Registry` during the FSM's initialization, simplifying the user's workflow configuration for common patterns.

### 12.5. Observability
The observability features include:
- Structured logging using the `log/slog` package
- Prometheus metrics for monitoring transition counts and durations
- OpenTelemetry tracing for distributed tracing

---

## 13. API Design Principles

### 13.1. Context Propagation
All operations accept a `context.Context` parameter to support:
- Cancellation
- Timeouts
- Distributed tracing

### 13.2. Error Handling
Errors are handled using Go's idiomatic approach:
- Functions return an error value as the last return parameter
- Errors are wrapped using `fmt.Errorf` with `%w` verb for error inspection
- Sentinel errors are provided for common failure cases
- Successful operations return a `TransitionResult` struct containing the new state, any auto-event, and persistence data

### 13.3. Concurrency Safety
The library is designed to be safe for concurrent use:
- The registry uses mutexes to protect shared state
- The state machine engine is stateless (state is passed as parameters)
- All public APIs are safe for concurrent access
- The `TransitionResult` struct is immutable after creation

### 13.4. Extensibility
The library is designed to be extensible:
- Users provide their own implementations of conditions and actions
- The registry allows for dynamic registration of implementations
- Observability features can be customized or extended
- The `TransitionResult` struct can be extended in future versions without breaking existing code

---

## 14. Performance Considerations

### 14.1. Memory Efficiency
- States and transitions are stored in optimized data structures
- Minimal memory allocation during transitions
- Efficient lookup of states and transitions

### 14.2. Execution Speed
- Conditions and actions are executed in a predictable order
- Minimal overhead for state transitions
- Benchmarks are provided to measure performance

### 14.3. Scalability
- The library is designed to handle high-throughput scenarios
- Concurrency-safe design allows for parallel processing
- Resource usage is minimized to support large-scale deployments

---

## 15. Testing Strategy

### 15.1. Unit Tests
- Each component has comprehensive unit tests
- Edge cases and error conditions are thoroughly tested
- Test coverage is measured and maintained

### 15.2. Integration Tests
- End-to-end workflows are tested
- Integration with external systems is verified
- Performance benchmarks are included

### 15.3. Example Workflows
- Real-world example workflows are provided
- Examples demonstrate best practices
- Examples serve as documentation and tests

---

## 16. Documentation

### 16.1. Godoc Comments
- All public APIs are documented with godoc comments
- Examples are provided for complex APIs
- Documentation is kept up-to-date with code changes

### 16.2. README
- Project overview and quick start guide
- Installation and usage instructions
- Contribution guidelines

### 16.3. Examples
- Runnable examples demonstrating common use cases
- Examples for advanced features and patterns
- Examples serve as both documentation and tests

---

## 17. Release Process

### 17.1. Versioning
- Semantic versioning is used (MAJOR.MINOR.PATCH)
- Breaking changes result in MAJOR version increments
- New features result in MINOR version increments
- Bug fixes result in PATCH version increments

### 17.2. Release Checklist
- All tests pass
- Documentation is up-to-date
- Examples are working
- CHANGELOG is updated
- Version is tagged in Git

---

## 18. Future Enhancements

### 18.1. Persistence
- Support for persisting state machine state
- Integration with popular databases
- Recovery from failures

### 18.2. Visualization
- Tools for visualizing workflow definitions
- Real-time monitoring dashboards
- Debugging and tracing tools

### 18.3. Advanced Features
- Support for parallel transitions
- Hierarchical state machines
- Time-based transitions