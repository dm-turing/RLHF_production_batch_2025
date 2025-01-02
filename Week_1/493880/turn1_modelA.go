package main

import (
	"fmt"
	"sync"
	"time"
)

// Subject represents an object with observers
type Subject interface {
	RegisterObserver(Observer)
	DeregisterObserver(Observer)
	NotifyObservers()
}

// ConcreteSubject is a specific subject that manages observers
type ConcreteSubject struct {
	observers []Observer
	mutex     sync.Mutex
}

func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{}
}

func (cs *ConcreteSubject) RegisterObserver(o Observer) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.observers = append(cs.observers, o)
}

func (cs *ConcreteSubject) DeregisterObserver(o Observer) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for i, observer := range cs.observers {
		if observer == o {
			cs.observers = append(cs.observers[:i], cs.observers[i+1:]...)
			break
		}
	}
}

func (cs *ConcreteSubject) NotifyObservers() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for _, observer := range cs.observers {
		observer.Update()
	}
}

// Observer represents an object that is observed
type Observer interface {
	Update()
}

// ConcreteObserver is a specific observer
type ConcreteObserver struct {
	name string
}

func NewConcreteObserver(name string) *ConcreteObserver {
	return &ConcreteObserver{name}
}

func (co *ConcreteObserver) Update() {
	fmt.Println(co.name, "is notified")
}

// Worker is a concurrent worker that handles notifications
type Worker struct {
	id      string
	subject Subject
	quit    chan struct{}
}

func NewWorker(id string, subject Subject) *Worker {
	return &Worker{
		id:      id,
		subject: subject,
		quit:    make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.quit:
				fmt.Println(w.id, "is quitting")
				return
			default:
				w.subject.NotifyObservers()
				fmt.Println(w.id, "notified observers")
				select {
				case <-time.After(1 * time.Second):
					continue
				case <-w.quit:
					fmt.Println(w.id, "is quitting")
					return
				}
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quit)
}

func main() {
	subject := NewConcreteSubject()
	observer1 := NewConcreteObserver("Observer 1")
	observer2 := NewConcreteObserver("Observer 2")

	subject.RegisterObserver(observer1)
	subject.RegisterObserver(observer2)

	worker := NewWorker("Worker", subject)
	worker.Start()

	time.Sleep(3 * time.Second)

	worker.Stop()
}
