package main

import (
	"context"
	"errors"
	"log"
)

// MockSecretManager is a mock implementation of a secret management service.
type MockSecretManager struct {
	credentials map[string]string
}

func NewMockSecretManager() *MockSecretManager {
	return &MockSecretManager{
		credentials: map[string]string{
			"user:alice": "secretTokenAlice",
		},
	}
}

// GetSecret retrieves a secret from the mock secret manager.
func (sm *MockSecretManager) GetSecret(ctx context.Context, key string) (string, error) {
	return sm.credentials[key], nil
}

// ValidateCredentials validates the retrieved credentials.
func ValidateCredentials(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("empty token")
	}
	return nil
}

// ProcessRequest uses the validated credentials to process a request.
func ProcessRequest(ctx context.Context, token string) error {
	log.Printf("Processing request with token: %s\n", token)
	return nil
}

// Callback represents a function that can be chained.
type Callback func(context.Context) error

// RunChainedCallbacks runs a sequence of callbacks.
func RunChainedCallbacks(ctx context.Context, callbacks ...Callback) error {
	for _, cb := range callbacks {
		if err := cb(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RetrieveCredentialsFromSecretManager is a callback function to retrieve credentials.
func RetrieveCredentialsFromSecretManager(ctx context.Context, secretManager *MockSecretManager, key string) (string, error) {
	token, err := secretManager.GetSecret(ctx, key)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateCallbackChain creates a callback chain for managing credentials.
func CreateCallbackChain(ctx context.Context, secretManager *MockSecretManager, key string) error {
	token, err := RetrieveCredentialsFromSecretManager(ctx, secretManager, key)
	if err != nil {
		return err
	}

	callbacks := []Callback{
		func(ctx context.Context) error {
			return ValidateCredentials(ctx, token)
		},
		func(ctx context.Context) error {
			return ProcessRequest(ctx, token)
		},
	}

	return RunChainedCallbacks(ctx, callbacks...)
}

func main() {
	ctx := context.Background()
	secretManager := NewMockSecretManager()

	if err := CreateCallbackChain(ctx, secretManager, "user:alice"); err != nil {
		log.Fatalf("Failed to process credentials: %v\n", err)
	}

	log.Println("Credentials processed successfully.")
}
