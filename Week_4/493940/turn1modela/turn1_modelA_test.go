package turn1modela

import "testing"

func TestUserDatabase(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer func() {
		if err := teardownDatabase(db); err != nil {
			t.Fatalf("Failed to teardown database: %v", err)
		}
	}()

	// Populate and verify user data
	expectedUsers, err := populateDatabase(db, 100)
	if err != nil {
		t.Fatalf("Failed to populate database: %v", err)
	}

	actualUsers, err := getUsersFromDatabase(db)
	if err != nil {
		t.Fatalf("Failed to get users from database: %v", err)
	}

	if len(actualUsers) != len(expectedUsers) {
		t.Fatalf("Expected %d users, got %d", len(expectedUsers), len(actualUsers))
	}

	for i, expected := range expectedUsers {
		actual := actualUsers[i]
		if actual.Name != expected.Name || actual.Email != expected.Email {
			t.Errorf("User mismatch: expected %+v, got %+v", expected, actual)
		}
	}
}
