package core

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// LogStartAction logs the start of the workflow
func LogStartAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Starting workflow...")
	return nil, nil
}

// LogProcessingAction logs the processing step
func LogProcessingAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Processing order...")
	return nil, nil
}

// LogCompletionAction logs the completion of the workflow
func LogCompletionAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Order completed successfully")
	return nil, nil
}

// LogFailureAction logs a failure
func LogFailureAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Order processing failed")
	return nil, nil
}

// IsUserValidCondition checks if the user is valid
func IsUserValidCondition(ctx context.Context, data map[string]any) (bool, error) {
	// Simulate user validation
	time.Sleep(100 * time.Millisecond)
	return true, nil // Always valid for this example
}

// IsPaymentSuccessCondition checks if payment was successful
func IsPaymentSuccessCondition(ctx context.Context, data map[string]any) (bool, error) {
	// Simulate payment processing with some randomness
	time.Sleep(200 * time.Millisecond)
	return rand.Float32() > 0.2, nil // 80% success rate
}

// ChargePaymentAction charges the user's payment method
func ChargePaymentAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Charging payment...")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("Payment charged successfully")

	// Return updated data
	return map[string]any{
		"paymentStatus": "charged",
		"chargedAt":     time.Now(),
	}, nil
}

// SendReceiptAction sends a receipt to the user
func SendReceiptAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Sending receipt...")
	time.Sleep(150 * time.Millisecond)
	fmt.Println("Receipt sent successfully")

	return map[string]any{
		"receiptSent": true,
		"sentAt":      time.Now(),
	}, nil
}

// HandleFailureAction handles a failure
func HandleFailureAction(ctx context.Context, data map[string]any) (map[string]any, error) {
	fmt.Println("Handling failure...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Failure handled")

	return map[string]any{
		"failureHandled": true,
		"handledAt":      time.Now(),
	}, nil
}
