package main

import (
	"math/rand"
	"net/http"
	"time"
)

// ABTestManager interface
type ABTestManager interface {
	GetVersion(r *http.Request) string
	AssignUser(r *http.Request) (string, error)
	SetVersionCookie(w http.ResponseWriter, version string)
	IsUserAssigned(r *http.Request, version string) bool
	GetWeightedVersion(weights map[string]int) string
	GetGeographicallyTargetedVersion(r *http.Request, versions map[string][]string) string
}

// BasicABTestManager is a simple implementation
type BasicABTestManager struct {
}

// ExtendedABTestManager uses cookies for persistent assignments
type ExtendedABTestManager struct {
	cookieName string
}

func NewExtendedABTestManager(cookieName string) *ExtendedABTestManager {
	return &ExtendedABTestManager{cookieName: cookieName}
}

// ... (other methods implemented as before)

// SetVersionCookie sets a cookie with the provided version
func (m *ExtendedABTestManager) SetVersionCookie(w http.ResponseWriter, version string) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    version,
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // Set the cookie to last for one week
		HttpOnly: true,
	})
}

// IsUserAssigned checks if the user is already assigned to the specified version
func (m *ExtendedABTestManager) IsUserAssigned(r *http.Request, version string) bool {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return false
	}
	return cookie.Value == version
}

// GetWeightedVersion assigns users based on a provided map of weights
func (m *ExtendedABTestManager) GetWeightedVersion(weights map[string]int) string {
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}
	if totalWeight == 0 {
		return "" // Return empty string if no weights are provided
	}

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(totalWeight)

	currentWeight := 0
	for version, weight := range weights {
		currentWeight += weight
		if randomNum < currentWeight {
			return version
		}
	}

	return "" // This should never be reached, but returning empty string just in case
}

// GetGeographicallyTargetedVersion assigns users based on a map of versions targeted to specific countries
func (m *ExtendedABTestManager) GetGeographicallyTargetedVersion(r *http.Request, versions map[string][]string) string {
	country := getCountryFromIP(r.RemoteAddr) // Implement getCountryFromIP function based on IP address
	for _, version := range versions[country] {
		if m.IsUserAssigned(r, version) {
			return version
		}
	}

	// Assign user to a random version within the targeted countries
	var targetVersions []string
	for _, versions := range versions {
		targetVersions = append(targetVersions, versions...)
	}
	if len(targetVersions) == 0 {
		return "" // No targeted versions found
	}
	return targetVersions[rand.Intn(len(targetVersions))]
}
