package main

import (
	"fmt"
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

func main() {
	// Create the application state
	state := NewState()

	// Create and register a callback function to handle UI updates based on the state
	state.Subscribe(func(newState *State) {
		// Abstracted GUI update logic
		fmt.Println("--- To-Do List ---")
		for _, todo := range newState.GetTodos() {
			status := "Done"
			if !todo.Done {
				status = "Pending"
			}
			fmt.Printf("[%d] %s - %s\n", todo.Index+1, todo.Text, status)
		}
		fmt.Println("-------------------")
	})

	// Application logic starts here
	state.AddTodo("Buy milk")
	state.AddTodo("Learn Go")
	state.ToggleTodo(0)
	state.RemoveTodo(1)
}
