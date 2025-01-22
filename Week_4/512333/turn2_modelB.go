package main

import (
	"fmt"
	"sync"
)

// TodoItem represents a single to-do item in the list.
type TodoItem struct {
	Text  string
	Done  bool
	Index int
}

// State represents the entire application state.
type State struct {
	todos     []TodoItem
	callbacks []func(newState *State)
}

// NewState creates a new instance of the application state.
func NewState() *State {
	return &State{}
}

// Subscribe registers a callback function to be called when the state changes.
func (s *State) Subscribe(cb func(newState *State)) {
	s.callbacks = append(s.callbacks, cb)
}

// NotifyAll calls all registered callbacks with the updated state.
func (s *State) NotifyAll() {
	for _, cb := range s.callbacks {
		cb(s)
	}
}

// AddTodo adds a new to-do item to the list and notifies all callbacks.
func (s *State) AddTodo(text string) {
	newItem := TodoItem{Text: text, Done: false, Index: len(s.todos)}
	s.todos = append(s.todos, newItem)
	s.NotifyAll()
}

// ToggleTodo toggles the done status of the to-do item at the specified index and notifies all callbacks.
func (s *State) ToggleTodo(index int) {
	if index >= 0 && index < len(s.todos) {
		s.todos[index].Done = !s.todos[index].Done
		s.NotifyAll()
	}
}

// RemoveTodo removes the to-do item at the specified index and notifies all callbacks.
func (s *State) RemoveTodo(index int) {
	if index >= 0 && index < len(s.todos) {
		s.todos = append(s.todos[:index], s.todos[index+1:]...)
		for i := index; i < len(s.todos); i++ {
			s.todos[i].Index = i
		}
		s.NotifyAll()
	}
}

// GetTodos returns the current list of to-do items.
func (s *State) GetTodos() []TodoItem {
	return s.todos
}

// messageBus is a simple distributed message bus
type messageBus struct {
	ch          chan interface{}
	subscribers []chan interface{}
	mtx         sync.Mutex
}

func newMessageBus() *messageBus {
	return &messageBus{
		ch:          make(chan interface{}),
		subscribers: make([]chan interface{}, 0),
	}
}

func (mb *messageBus) publish(msg interface{}) {
	mb.ch <- msg
}

func (mb *messageBus) subscribe() <-chan interface{} {
	ch := make(chan interface{})
	mb.mtx.Lock()
	defer mb.mtx.Unlock()
	mb.subscribers = append(mb.subscribers, ch)
	return ch
}

func (mb *messageBus) unsubscribe(ch chan interface{}) {
	mb.mtx.Lock()
	defer mb.mtx.Unlock()
	var subscribers []chan interface{}
	for _, sub := range mb.subscribers {
		if sub != ch {
			subscribers = append(subscribers, sub)
		}
	}
	mb.subscribers = subscribers
}

func (mb *messageBus) start() {
	go func() {
		for msg := range mb.ch {
			mb.mtx.Lock()
			for _, sub := range mb.subscribers {
				sub <- msg
			}
			mb.mtx.Unlock()
		}
	}()
}

// Main function becomes distrubuted
func main() {
	// Create the application state and message bus
	state := NewState()
	mb := newMessageBus()
	mb.start()

	// Modify the addTodo method to publish state updates to the message bus.
	state.AddTodo = func(text string) {
		state.Mutex.Lock()
		defer state.Mutex.Unlock()
		newItem := TodoItem{Text: text, Done: false, Index: len(state.todos)}
		state.todos = append(state.todos, newItem)
		mb.publish(state) // Notify all subscribers through the message bus
	}

	// Subscribe a different process to handle updates
	go func() {
		ch := mb.subscribe()
		defer mb.unsubscribe(ch)
		for newState := range ch {
			// Another process can now apply state updates appropriately
			nState := newState.(*State)
			fmt.Println("\n**Updating other process state:**")
			for _, todo := range nState.GetTodos() {
				fmt.Printf("[%d] %s - %t\n", todo.Index+1, todo.Text, todo.Done)
			}
		}
	}()

	// Application logic starts here, as before
	state.AddTodo("Task 1")
	state.AddTodo("Task 2")
	state.ToggleTodo(1)
	state.RemoveTodo(0)
}
