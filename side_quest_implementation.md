# GoMachina: Side Quest Implementation Details

## 1. Introduction

This document provides a detailed design and implementation guide for handling "Side Quests" within the `GoMachina` state machine framework. Side quests represent temporary diversions from a main workflow, allowing users to complete optional tasks or sub-processes before returning to their original position.

## 2. The "Side Quest" Problem Revisited

In complex user journeys (e.g., `A -> B -> C -> D -> E`), a user might need to temporarily jump to a "side quest" state, like `B#`, from a main flow state like `D`. After completing `B#`, they must seamlessly return to `D` to continue.

Key challenges addressed by this design include:
*   **Remembering Return Points:** How does the system know where to return after the side quest?
*   **Handling Drop-offs:** What happens if a user abandons the application while in a side quest? Should they resume the side quest or be reverted to the main flow?
*   **UI Synchronization:** Ensuring the UI accurately reflects the user's position, even during temporary diversions.

## 3. Core Concepts for Side Quests

### 3.1. Workflow Stack

The `WorkflowStack` is a list of state names stored within the `Context` object. It acts like a call stack in programming, remembering the sequence of states the user has temporarily diverted from.

*   **Purpose:** To store the "return address" when entering a side quest.
*   **Mechanism:** When entering a side quest, the current state is "pushed" onto the stack. When exiting, the last state is "popped" from the stack.

### 3.2. `isSideQuest` Flag

A boolean flag added to the `State` definition in `workflow.yaml` to explicitly mark a state as a side quest.

*   **Purpose:** Allows the application layer to easily identify side quest states for specific business rules (e.g., drop-off handling).
*   **Mechanism:** `isSideQuest: true` in the YAML.

### 3.3. Dynamic Transition Target (`__next_state_override`)

A special mechanism where an `ActionFunc` can return a specific key (`__next_state_override`) in its `persistenceData` to dynamically determine the next state, overriding the static `target` in the YAML.

*   **Purpose:** Enables transitions to a state that is not known statically, such as returning to a state popped from the `WorkflowStack`.

### 3.4. Pre-defined Action: `__RETURN_TO_PREVIOUS_STATE__`

A special `ActionFunc` internally registered by `GoMachina` for convenience.

*   **Purpose:** To simplify the common operation of popping a state from the `WorkflowStack` and signaling it as the dynamic return target.
*   **Mechanism:** When this action is executed, it pops the top state from the `WorkflowStack` and returns it as the `__next_state_override`.

## 4. Detailed Implementation

### 4.1. `workflow.yaml` Configuration

This example shows how to define a side quest (`B#`) accessible from state `D`, and how to return from it.

```yaml
states:
  D:
    # ... other transitions for D ...
    transitions:
      - event: "start_b_sharp_quest" # User clicks to start side quest
        target: "B#"                 # Actual state B#
        actions:
          - "pushCurrentStateToStack" # Action to push "D" onto the stack

  B#:
    isSideQuest: true # Mark this state as a side quest
    onEnter:
      - "logEnteringBSharp" # Log entry into the side quest
    onLeave:
      - "logLeavingBSharp"  # Log exit from the side quest
    transitions:
      - event: "submit_b_sharp_form" # User submits form within B#
        target: "B#"                 # Stays in B# for further processing
        actions:
          - "processBSharpForm"
      - event: "confirm_b_sharp_and_return" # User confirms completion and wants to return
        target: "" # Dynamic target, determined by action
        actions:
          - "saveBSharpResults" # Save final data from B#
          - "__RETURN_TO_PREVIOUS_STATE__" # Pre-defined action to pop stack and return
```

### 4.2. Go `ActionFunc` Implementations

These are the conceptual Go functions that implement the side quest logic.

```go
package main // Or your internal/core/actions package

import (
	"context"
	"fmt"
	"github.com/your-username/go-machina/machina" // Assuming your module path
)

// pushCurrentStateToStack: Pushes the state from which the side quest was initiated onto the WorkflowStack.
// This action is called on the transition *into* the side quest (e.g., D -> B#).
func pushCurrentStateToStack(ctx machina.Context, payload interface{}) (map[string]interface{}, error) {
    // The 'currentState' is passed as an argument to fsm.Trigger().
    // We need to get it from the Context if it's not directly available in payload.
    // For simplicity, let's assume the application passes it in payload for this action.
    sourceState := payload.(map[string]interface{})["source_state"].(string) // e.g., "D"

    // Access the WorkflowStack from the Context
    // (Assuming Context is designed to hold and allow modification of WorkflowStack)
    var currentStack []string
    if stack, ok := ctx.Value("workflow_stack").([]string); ok {
        currentStack = stack
    }
    
    newStack := append(currentStack, sourceState)
    fmt.Printf("Pushing %s onto stack. New stack: %v\n", sourceState, newStack)

    // Return the updated stack for the application to persist.
    return map[string]interface{}{"workflow_stack": newStack}, nil
}

// processBSharpForm: Processes data submitted within the B# side quest.
func processBSharpForm(ctx machina.Context, payload interface{}) (map[string]interface{}, error) {
    formData := payload.(map[string]interface{})["b_sharp_data"].(string)
    fmt.Printf("Processing B# form data: %s\n", formData)
    // ... perform validation, call external services, etc. ...
    return map[string]interface{}{"b_sharp_processed": true}, nil
}

// saveBSharpResults: Saves final data from B# before returning.
func saveBSharpResults(ctx machina.Context, payload interface{}) (map[string]interface{}, error) {
    fmt.Println("Saving final B# results.")
    // ... persist data to database ...
    return nil, nil
}

// Note: The __RETURN_TO_PREVIOUS_STATE__ ActionFunc is implemented internally by GoMachina.
// It pops from the WorkflowStack and signals the dynamic target.
```

### 4.3. Application Layer Logic (Drop-off Handling)

This logic resides in your application's API handler that serves the "resume journey" endpoint (e.g., `GET /api/status`).

```go
package main // Or your internal/api package

import (
	"context"
	"fmt"
	"time"
	"github.com/your-username/go-machina/machina" // Assuming your module path
)

// Assume these are loaded at app startup
var fsm *machina.FSM
var workflowDefinition *machina.Definition
// Assume dataService is available to load/save user data and context

func handleResumeJourney(w http.ResponseWriter, r *http.Request) {
    userID := "user123" // Get user ID from session/auth
    
    // 1. Load user's current state and context from DB
    userState, userContextData, err := dataService.LoadUserWorkflowData(userID)
    if err != nil { /* handle error */ return } 

    // 2. Check if current state is a side quest
    stateDef, ok := workflowDefinition.States[userState]
    if !ok { /* handle error */ return } 

    if stateDef.IsSideQuest { 
        // 3. Check for drop-off condition (e.g., inactivity)
        lastActivityTime := userContextData["last_activity_time"].(time.Time) // Assuming this is tracked
        if time.Since(lastActivityTime) > 30*time.Minute { // 30 min inactivity
            fmt.Printf("User %s dropped off in side quest %s. Forcing return.\n", userID, userState)
            
            // 4. Force return to previous state
            // The 'force_return_event' transition in workflow.yaml for B# would use __RETURN_TO_PREVIOUS_STATE__ action.
            // The payload needs to contain the current WorkflowStack for the action to pop from.
            payload := map[string]interface{}{
                "workflow_stack": userContextData["workflow_stack"], // Pass the stack to the action
            }
            
            // Call Trigger with the internal force event
            newState, persistenceData, triggerErr := fsm.Trigger(
                r.Context(), // Use request context
                userState,
                "force_return_from_side_quest",
                payload,
            )
            if triggerErr != nil { /* handle trigger error */ return }

            // 5. Persist updated state and context
            userContextData["workflow_stack"] = persistenceData["workflow_stack"] // Update stack
            dataService.SaveUserWorkflowData(userID, newState, userContextData)

            // 6. Return directives for the new state (main flow)
            returnDirectives(w, workflowDefinition.States[newState].ClientDirectives)
            return
        }
    }
    
    // 7. If not a side quest or no drop-off, return directives for current state
    returnDirectives(w, stateDef.ClientDirectives)
}

func returnDirectives(w http.ResponseWriter, directives map[string]interface{}) {
    // ... marshal directives to JSON and write to response ...
}

---
