package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
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

func init() {
	registry.Register("Upper", func(s string) (string, error) {
		return strings.ToUpper(s), nil
	})
	registry.Register("Lower", func(s string) (string, error) {
		return strings.ToLower(s), nil
	})
	registry.Register("TrimSpace", func(s string) (string, error) {
		return strings.TrimSpace(s), nil
	})
	registry.Register("Contains", func(s string) (string, error) {
		subStr := strings.Fields(s)[1]
		inputStr := strings.Fields(s)[0]
		contains := strings.Contains(inputStr, subStr)
		var result string
		if contains {
			result = "true"
		} else {
			result = "false"
		}
		return result, nil
	})
	registry.Register("Index", func(s string) (string, error) {
		subStr := strings.Fields(s)[1]
		inputStr := strings.Fields(s)[0]
		index := strings.Index(inputStr, subStr)
		return fmt.Sprintf("%d", index), nil
	})
	registry.Register("Join", func(s string) (string, error) {
		separator := strings.Fields(s)[0]
		args := strings.Fields(s)[1:]
		return strings.Join(args, separator), nil
	})
	registry.Register("Split", func(s string) (string, error) {
		separator := strings.Fields(s)[1]
		inputStr := strings.Fields(s)[0]
		result := strings.Split(inputStr, separator)
		return strings.Join(result, " "), nil
	})
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the Go strings library REPL!")
	fmt.Println("Available functions: Upper, Lower, TrimSpace, Contains, Index, Join, Split")
	fmt.Println("Enter a command or 'exit' to quit.")
	for {
		fmt.Print("→ ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			break
		}
		cmdArgs := strings.Fields(input)
		cmd := cmdArgs[0]
		args := strings.Join(cmdArgs[1:], " ")
		result, err := registry.Invoke(cmd, args)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
	fmt.Println("Goodbye!")
}
