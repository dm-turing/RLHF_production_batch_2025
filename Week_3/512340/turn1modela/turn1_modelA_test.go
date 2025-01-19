// auth_test.go
package auth

import (
	"testing"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	store := NewUserStore()

	// Test Registration Success
	err := store.Register("testuser", "StrongPass1!")
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}

	// Test Authentication Success
	if !store.Authenticate("testuser", "StrongPass1!") {
		t.Error("authentication failed; user should be able to log in with correct credentials")
	}

	// Test Registration Failure (username already exists)
	if err := store.Register("testuser", "AnotherPass1!"); err == nil {
		t.Error("expected registration failure; username should already exist")
	}

	// Test Registration Failure (weak password)
	if err := store.Register("newuser", "weak"); err == nil {
		t.Error("expected registration failure; password too weak")
	}
}

func TestPasswordRequirements(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"noNumbersOrSpecial!", false},
		{"NoSpecialChar123", false},
		{"Sh0rt!", false},                // too short
		{"ValidPass123!", true},          // valid
		{"another#Valid1", true},         // valid
		{"WayTo0shor", false},            // no special
	}

	for _, test := range tests {
		if valid := isValidPassword(test.password); valid != test.valid {
			t.Errorf("expected password validity for %q to be %v, got %v", test.password, test.valid, valid)
		}
	}
}
