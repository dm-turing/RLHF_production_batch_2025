package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

// ABTestManager interface
type ABTestManager interface {
	GetVersion(r *http.Request) string
	AssignUser(r *http.Request) (string, error)
	SetWeights(weights map[string]int)
	GetWeights() map[string]int
	UpdateVersion(version string) error
}

// BasicABTestManager is a simple implementation
type BasicABTestManager struct {
	weights map[string]int
	version string
}

func (m *BasicABTestManager) GetVersion(r *http.Request) string {
	return m.version
}

func (m *BasicABTestManager) AssignUser(r *http.Request) (string, error) {
	versions := make([]string, 0, len(m.weights))
	for version, weight := range m.weights {
		for i := 0; i < weight; i++ {
			versions = append(versions, version)
		}
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions available")
	}
	return versions[rand.Intn(len(versions))], nil
}

func (m *BasicABTestManager) SetWeights(weights map[string]int) {
	m.weights = weights
}

func (m *BasicABTestManager) GetWeights() map[string]int {
	return m.weights
}

func (m *BasicABTestManager) UpdateVersion(version string) error {
	if _, ok := m.weights[version]; !ok {
		return fmt.Errorf("version %s not found", version)
	}
	m.version = version
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		manager := BasicABTestManager{
			weights: map[string]int{"A": 50, "B": 50}, // Default weights
			version: "A",
		}

		version, err := manager.AssignUser(r)
		if err != nil {
			fmt.Fprintf(w, "Error assigning user: %v\n", err)
			return
		}

		fmt.Fprintf(w, "You are assigned to version %s\n", version)
	})

	http.ListenAndServe(":8080", nil)
}
