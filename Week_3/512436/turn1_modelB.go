package main

import (
	"fmt"
)

// PaymentError is a custom error type to handle payment processing errors
type PaymentError struct {
	Code    string
	Message string
}

// Error implements the error interface for PaymentError
func (e PaymentError) Error() string {
	return fmt.Sprintf("Payment Error: [%s] %s", e.Code, e.Message)
}

// NewPaymentError creates a new PaymentError instance
func NewPaymentError(code, message string) PaymentError {
	return PaymentError{Code: code, Message: message}
}

// PaymentError codes
const (
	InvalidOrderErrorCode         = "INVALID_ORDER"
	InsufficientFundsErrorCode    = "INSUFFICIENT_FUNDS"
	CardDeclinedErrorCode         = "CARD_DECLINED"
	ProcessingErrorCode           = "PROCESSING_ERROR"
	DuplicateTransactionErrorCode = "DUPLICATE_TRANSACTION"
	AuthenticationErrorCode       = "AUTHENTICATION_ERROR"
)

// Example usage of the DSL
func main() {
	// Order placement error
	err := NewPaymentError(InvalidOrderErrorCode, "Invalid order details provided")
	handlePaymentError(err)

	// Insufficient funds error
	err = NewPaymentError(InsufficientFundsErrorCode, "Insufficient balance in the payment account")
	handlePaymentError(err)

	// Card declined error
	err = NewPaymentError(CardDeclinedErrorCode, "Card declined due to insufficient funds")
	handlePaymentError(err)

	// Processing error
	err = NewPaymentError(ProcessingErrorCode, "An internal processing error occurred")
	handlePaymentError(err)

	// Duplicate transaction error
	err = NewPaymentError(DuplicateTransactionErrorCode, "Duplicate transaction detected")
	handlePaymentError(err)

	// Authentication error
	err = NewPaymentError(AuthenticationErrorCode, "Authentication failed for the payment method")
	handlePaymentError(err)

	// No error (successful transaction)
	handlePaymentError(nil)
}

// Function to handle payment errors using the DSL
func handlePaymentError(err error) {
	switch err.(type) {
	case PaymentError:
		pe := err.(PaymentError)
		switch pe.Code {
		case InvalidOrderErrorCode:
			fmt.Println("Error: Invalid order -", pe.Message)
		case InsufficientFundsErrorCode:
			fmt.Println("Error: Insufficient funds -", pe.Message)
		case CardDeclinedErrorCode:
			fmt.Println("Error: Card declined -", pe.Message)
		case ProcessingErrorCode:
			fmt.Println("Error: Processing error -", pe.Message)
		case DuplicateTransactionErrorCode:
			fmt.Println("Error: Duplicate transaction -", pe.Message)
		case AuthenticationErrorCode:
			fmt.Println("Error: Authentication error -", pe.Message)
		default:
			fmt.Println("Unrecognized payment error:", pe.Error())
		}
	case nil:
		fmt.Println("Transaction successful!")
	default:
		fmt.Println("Unexpected error:", err)
	}
}
