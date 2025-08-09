# GoMachina

A declarative state machine engine for Go applications.

## Overview

GoMachina is a powerful, declarative, and observable state machine library for orchestrating complex business workflows with safety, clarity, and extensibility.

## Features

- **Declarative-First**: Define state machine structure in YAML/JSON configuration files
- **High Observability**: Built-in structured logging, metrics, and tracing
- **Safety & Predictability**: Safe execution order and concurrent use
- **Extensible & Decoupled**: Generic, domain-agnostic design with user-provided implementations

## Installation

```bash
go get github.com/rahulpahuja/go-machina
```

## Usage

1. Define your workflow in a YAML file:

```yaml
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

  complete:
    name: complete
    onEnter:
      - "sendReceipt"
```

2. Implement your conditions and actions using function types:

```go
// Condition function type
type ConditionFunc func(ctx context.Context, data map[string]any) (bool, error)

// Action function type
type ActionFunc func(ctx context.Context, data map[string]any) (map[string]any, error)

// Example condition implementation
func IsUserValidCondition(ctx context.Context, data map[string]any) (bool, error) {
    // Your validation logic here
    return true, nil
}

// Example action implementation
func ChargePaymentAction(ctx context.Context, data map[string]any) (map[string]any, error) {
    // Your payment logic here
    return map[string]any{
        "paymentStatus": "charged",
    }, nil
}

// Example onEnter action
func SendReceiptAction(ctx context.Context, data map[string]any) (map[string]any, error) {
    // Your receipt logic here
    return map[string]any{
        "receiptSent": true,
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
registry.RegisterAction("sendReceipt", SendReceiptAction)

// Create state machine with observability options
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
fsm := machina.NewStateMachine(definition, registry, logger,
    machina.WithMetrics(prometheus.DefaultRegisterer),
    machina.WithTracer(otel.Tracer("my-app")))

// Trigger events to process the workflow
ctx := context.Background()
data := map[string]any{
    "userId": "user123",
    "amount": 99.99,
}

// Trigger the validate event
newState, result, err := fsm.Trigger(ctx, "start", "validate", data)
if err != nil {
    log.Fatal(err)
}

// Trigger the process event
finalState, finalResult, err := fsm.Trigger(ctx, newState, "process", result)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Workflow completed. Final state: %s, Result: %v\n", finalState, finalResult)
```

## Running the Example

```bash
go run cmd/server/main.go
```

## Testing

```bash
go test ./machina
```

## License

MIT