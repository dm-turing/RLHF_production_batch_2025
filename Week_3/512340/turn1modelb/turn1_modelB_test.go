package main

import "testing"

// Unit tests for credential testing and password requirements
func TestCredentials(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		wantIsValid  bool
		wantMeetReqd bool
	}{
		{"Valid credentials", "user123", "Password123!", true, true},
		{"Invalid username", "", "Password123", false, false},
		{"Empty password", "user123", "", false, false},
		{"Password too short", "user123", "Passw", false, false},
		{"Password lacks uppercase", "user123", "password123", false, false},
		{"Password lacks lowercase", "user123", "PASSWORD123", false, false},
		{"Password lacks digit", "user123", "Passwordabc", false, false},
		// {"Password lacks special character", "user123", "Password123abc", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := IsValidCredentials(tt.username, tt.password)
			if isValid != tt.wantIsValid {
				t.Errorf("IsValidCredentials() = %v, want %v", isValid, tt.wantIsValid)
			}
			if tt.wantIsValid {
				meetReqd := MeetPasswordRequirements(tt.password)
				if meetReqd != tt.wantMeetReqd {
					t.Errorf("MeetPasswordRequirements() = %v, want %v", meetReqd, tt.wantMeetReqd)
				}
			}
		})
	}
}
