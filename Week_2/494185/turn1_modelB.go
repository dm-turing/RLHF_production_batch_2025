package main

import (
	"fmt"
	"net"
)

// Custom error type for network errors
type NetworkError struct {
	Op  string
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %s: %v", e.Op, e.Err)
}

// Function to simulate multiple network errors
func SimulateNetworkErrors() ([]error, error) {
	var errors []error
	// Simulate a DNS lookup error
	dnsErr := &net.DNSError{
		Name: "example.com",
		Err:  "server misbehaving",
	}
	errors = append(errors, &NetworkError{Op: "DNS lookup", Err: dnsErr})

	// Simulate a connection timeout error
	timeoutErr := &net.OpError{
		Op:   "dial",
		Net:  "tcp",
		Addr: &net.TCPAddr{IP: net.IPv4(8, 8, 8, 8), Port: 53},
		Err:  fmt.Errorf("%s", "timeout error"),
	}
	errors = append(errors, &NetworkError{Op: "connect", Err: timeoutErr})
	return errors, nil
}

func main() {
	errors, _ := SimulateNetworkErrors()
	// Compose the error type with nested fields of errors
	composedError := &NetworkError{Op: "multiple network operations", Err: fmt.Errorf("%v", errors)}
	fmt.Println(composedError)
}
