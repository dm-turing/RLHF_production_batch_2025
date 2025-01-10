package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/gookit/color"
	"github.com/jdkato/prose"
)

type Callback func(string) (string, error)

type CallbackRegistry struct {
	callbacks map[string]Callback
	mu        sync.Mutex
}

func (r *CallbackRegistry) Register(name string, cb Callback) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.callbacks[name] = cb
}

func (r *CallbackRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.callbacks, name)
}

func (r *CallbackRegistry) Invoke(name, arg string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cb, ok := r.callbacks[name]
	if !ok {
		return "", fmt.Errorf("no callback registered for '%s'", name)
	}
	return cb(arg)
}

var registry = &CallbackRegistry{
	callbacks: make(map[string]Callback),
}

var availableFunctions = []string{"Upper", "Lower", "TrimSpace", "Contains", "Index", "Join", "Split"}

func init() {
	// Initialize autocorrection dictionary
	suggest.LoadDefaultDictionaries()

	// Register callback functions
	registry.Register("Upper", func(s string) (string, error) {
		return strings.ToUpper(s), nil
	})
	// ... (Other callback registrations remain the same)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the Go strings library REPL with Autocorrect and Autocompletion!")
	fmt.Println("Available functions: Upper, Lower, TrimSpace, Contains, Index, Join, Split")
	fmt.Println("Enter a command or 'exit' to quit.")

	terminal.SetColorMode(terminal.ColorModeAuto)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	for {
		fmt.Print("→ ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		// Handle autocompletion
		input, isCompleted := autocompletInput(input)
		if isCompleted {
			fmt.Println(input)
		}

		// Handle autocorrection
		input = autocorrectInput(input)

		if input == "exit" {
			break
		}
		cmdArgs := strings.Fields(input)
		cmd := cmdArgs[0]
		args := strings.Join(cmdArgs[1:], " ")

		result, err := registry.Invoke(cmd, args)
		if err != nil {
			red.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
	fmt.Println("Goodbye!")
}

func autocorrectInput(input string) string {
	cmd := strings.Fields(input)[0]
	if _, ok := registry.callbacks[cmd]; !ok {
		suggestions := suggest.Suggestions(cmd, availableFunctions)
		if len(suggestions) > 0 {
			yellow.Println("Did you mean:", suggestions[0])
			input = strings.Replace(input, cmd, suggestions[0], 1)
		}
	}
	return input
}

func autocompletInput(input string) (string, bool) {
	cmd := strings.Fields(input)[0]
	if len(cmd) > 0 {
		for _, f := range availableFunctions {
			if strings.HasPrefix(f, cmd) {