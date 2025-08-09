package machina

import (
	"context"
	"errors"
	"time"
)

// MockCondition is a test condition implementation
func MockTrueCondition(ctx context.Context, data map[string]any) (bool, error) {
	return true, nil
}

func MockFalseCondition(ctx context.Context, data map[string]any) (bool, error) {
	return false, nil
}

func MockErrorCondition(ctx context.Context, data map[string]any) (bool, error) {
	return false, errors.New("condition error")
}

func MockSlowCondition(ctx context.Context, data map[string]any) (bool, error) {
	time.Sleep(300 * time.Millisecond)
	return true, nil
}

// MockGuardCondition is a test guard condition
func MockGuardCondition(ctx context.Context, data map[string]any) (bool, error) {
	return true, nil
}

// MockFailingGuardCondition is a test guard condition that fails
func MockFailingGuardCondition(ctx context.Context, data map[string]any) (bool, error) {
	return false, nil
}

// MockAction is a test action implementation
func MockNoOpAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	return nil, nil
}

func MockErrorAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	return nil, errors.New("action error")
}

func MockSlowAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	select {
	case <-time.After(300 * time.Millisecond):
		return nil, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func MockUpdateAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	return map[string]any{
		"updated": true,
	}, nil
}
