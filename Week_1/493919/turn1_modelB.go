package main

import (
	"context"
	"fmt"
)

// Define a context key type for security credentials
type securityCredentialsKey int

const (
	securityCredentialsKeyID securityCredentialsKey = iota
	securityCredentialsKeyToken
)

// Define a SecurityCredentials struct to hold credentials
type SecurityCredentials struct {
	UserID      string
	AccessToken string
}

// Function to set security credentials in the context
func WithSecurityCredentials(ctx context.Context, credentials *SecurityCredentials) context.Context {
	return context.WithValue(ctx, securityCredentialsKeyID, credentials.UserID)
}

// Function to retrieve security credentials from the context
func GetSecurityCredentials(ctx context.Context) *SecurityCredentials {
	userID, ok := ctx.Value(securityCredentialsKeyID).(string)
	if !ok {
		return nil
	}
	accessToken, ok := ctx.Value(securityCredentialsKeyToken).(string)
	if !ok {
		return nil
	}
	return &SecurityCredentials{
		UserID:      userID,
		AccessToken: accessToken,
	}
}

// Example usage of context propagation
func main() {
	credentials := &SecurityCredentials{
		UserID:      "example_user",
		AccessToken: "example_token",
	}
	// Create a new context with security credentials
	ctx := WithSecurityCredentials(context.Background(), credentials)
	// Perform operations using the context with credentials
	performOperations(ctx)
}

// Example function that performs operations with the context
func performOperations(ctx context.Context) {
	// Retrieve credentials from the context
	credentials := GetSecurityCredentials(ctx)
	if credentials == nil {
		fmt.Println("Security credentials not found in the context.")
		return
	}
	// Use the credentials to perform operations
	fmt.Println("User ID:", credentials.UserID)
	fmt.Println("Access Token:", credentials.AccessToken)
}
