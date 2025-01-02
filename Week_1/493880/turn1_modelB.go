package main

import (
	"fmt"
	"sync"
)

// Subject Interface to define an event dispatcher
type Subject interface {
	RegisterObserver(observer Observer)
	NotifyObservers(data interface{})
}

// Observer Interface to define an event listener
type Observer interface {
	Notify(data interface{})
}

// Notifier struct represents a subject in the observer pattern.
type Notifier struct {
	observers []Observer
	mu        sync.Mutex
}

// RegisterObserver adds a new observer to the subject
func (n *Notifier) RegisterObserver(observer Observer) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.observers = append(n.observers, observer)
}

// NotifyObservers notifies all observers about the change in data
func (n *Notifier) NotifyObservers(data interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, observer := range n.observers {
		observer.Notify(data)
	}
}

// PriceMonitor represents an observer that observes the change in price.
type PriceMonitor struct {
	price float64
	mu    sync.Mutex
	wg    sync.WaitGroup
}

// Notify updates the internal price on notification.
func (pm *PriceMonitor) Notify(data interface{}) {
	if price, ok := data.(float64); ok {
		pm.mu.Lock()
		defer pm.mu.Unlock()
		pm.price = price
		pm.wg.Done()
	}
}

func (pm *PriceMonitor) WaitForNotification() {
	pm.wg.Add(1)
	pm.wg.Wait()
}

func (pm *PriceMonitor) GetPrice() float64 {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.price
}

// Main function to demonstrate the use of observer pattern.
func main() {
	// Create the subject (price data publisher)
	priceSubject := new(Notifier)

	// Create observers (price monitors)
	pm1 := &PriceMonitor{}
	pm2 := &PriceMonitor{}

	// Register observers with the subject
	priceSubject.RegisterObserver(pm1)
	priceSubject.RegisterObserver(pm2)

	// Simulate price updates in concurrent goroutines
	go func() {
		for i := 1; i <= 5; i++ {
			priceSubject.NotifyObservers(float64(i) * 100)
		}
	}()

	fmt.Println("Waiting for price notifications...")

	pm1.WaitForNotification()
	pm2.WaitForNotification()

	fmt.Println("Observer 1's Price:", pm1.GetPrice())
	fmt.Println("Observer 2's Price:", pm2.GetPrice())
}
