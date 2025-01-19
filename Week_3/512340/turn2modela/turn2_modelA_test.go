// auth_test.go
package auth

import (
	"testing"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	store := NewUserStore(defaultCredentialsValidation, defaultPasswordValidation)

	// Test Registration Success
	err := store.Register("testuser@example.com", "StrongPass1!")
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}

	// Test Authentication Success
	if !store.Authenticate("testuser@example.com", "StrongPass1!") {
		t.Error("authentication failed; user should be able to log in with correct credentials")
	}

	// Test Registration Failure (username already exists)
	if err := store.Register("testuser@example.com", "AnotherPass1!"); err == nil {
		t.Error("expected registration failure; username should already exist")
	}

	// Test Registration Failure (weak password)
	if err := store.Register("newuser@example.com", "weak"); err == nil {
		t.Error("expected registration failure; password too weak")
	}

	// Test Registration Failure (invalid email)
	if err := store.Register("invalidemail", "ValidPass1!"); err == nil {
		t.Error("expected registration failure; invalid email format")
	}
}

func TestPasswordRequirements(t *testing.T) {
	passwordTests := []struct {
		password string
		valid    bool
	}{
		{"noNumbersOrSpecial!", false},
		{"NoSpecialChar123", false},
		{"Sh0rt!", false},        // too short
		{"ValidPass123!", true},  // valid
		{"another#Valid1", true}, // valid
		{"WayTo0shor", false},    // no special
	}

	for _, test := range passwordTests {
		err := defaultPasswordValidation(test.password)
		valid := (err == nil)
		if valid != test.valid {
			t.Errorf("expected password validity for %q to be %v, got %v", test.password, test.valid, valid)
		}
	}
}

func TestCredentialRequirements(t *testing.T) {
	credentialTests := []struct {
		username string
		valid    bool
	}{
		{"valid@example.com", true},
		{"invalid-email", false},
		{"another.valid@domain.com", true},
		{"missingatsign.com", false},
	}

	for _, test := range credentialTests {
		err := defaultCredentialsValidation(test.username, "")
		valid := (err == nil)
		if valid != test.valid {
			t.Errorf("expected credential validity for %q to be %v, got %v", test.username, test.valid, valid)
		}
	}
}
