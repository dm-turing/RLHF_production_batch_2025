package main

import (
	"fmt"
)

type PaymentError interface {
	error
	// fmt.Stringer
}

type OrderPlacementError struct {
	code    string
	message string
}

func (e OrderPlacementError) Error() string {
	return fmt.Sprintf("Order Placement Error: %s, Code: %s", e.message, e.code)
}

type OrderConfirmationError struct {
	code    string
	message string
}

func (e OrderConfirmationError) Error() string {
	return fmt.Sprintf("Order Confirmation Error: %s, Code: %s", e.message, e.code)
}

type TransactionSuccessError struct {
	code    string
	message string
}

func (e TransactionSuccessError) Error() string {
	return fmt.Sprintf("Transaction Success Error: %s, Code: %s", e.message, e.code)
}

func newOrderPlacementError(code, message string) PaymentError {
	return &OrderPlacementError{code: code, message: message}
}

func newOrderConfirmationError(code, message string) PaymentError {
	return &OrderConfirmationError{code: code, message: message}
}

func newTransactionSuccessError(code, message string) PaymentError {
	return &TransactionSuccessError{code: code, message: message}
}

func handlePayment(process func() PaymentError) error {
	err := process()
	if err != nil {
		fmt.Println("Payment Processing Error:", err)
		return err
	}
	return nil
}

func processOrderPlacement() PaymentError {
	return newOrderPlacementError("O123", "Insufficient stock")
}

func processOrderConfirmation() PaymentError {
	return newOrderConfirmationError("C456", "Billing address mismatch")
}

func processTransaction() PaymentError {
	return newTransactionSuccessError("S789", "Payment processed successfully")
}

func main() {
	// Simulate payment processing
	handlePayment(processOrderPlacement)
	handlePayment(processOrderConfirmation)
	handlePayment(processTransaction)
}
