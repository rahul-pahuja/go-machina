package machina

import "context"

// ConditionFunc defines the function signature for evaluating transition conditions
type ConditionFunc func(ctx context.Context, data map[string]any) (bool, error)

// ActionFunc defines the function signature for executing state actions
// It returns a map of updated data and an error
type ActionFunc func(ctx context.Context, data map[string]any) (map[string]any, error)
