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

var words []string

func loadDictionary() {
	// Load a simple dictionary
	words = []string{
		"hello", "world", "to", "upper", "lower", "trimspace", "contains",
		"index", "join", "split", "exit",
	}
}

func correctWord(word string) string {
	minDist := 1000
	correctedWord := word
	for _, w := range words {
		dist := levenshteinDistance(word, w)
		if dist < minDist {
			minDist = dist
			correctedWord = w
		}
	}
	return correctedWord
}

func levenshteinDistance(s1, s2 string) int {
	lenS1, lenS2 := len(s1), len(s2)
	dp := make([][]int, lenS1+1)
	for i := 0; i <= lenS1; i++ {
		dp[i] = make([]int, lenS2+1)
		for j := 0; j <= lenS2; j++ {
			if i == 0 {
				dp[i][j] = j
			} else if j == 0 {
				dp[i][j] = i
			} else if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], // deletion
					dp[i][j-1],   // insertion
					dp[i-1][j-1]) // substitution
			}
		}
	}
	return dp[lenS1][lenS2]
}

func autoComplete(prefix string) []string {
	var completions []string
	for _, word := range words {
		if strings.HasPrefix(word, prefix) {
			completions = append(completions, word)
		}
	}
	return completions
}

func main() {
	loadDictionary()
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
		cmdArgs := strings.Fields(input)
		cmd := correctWord(cmdArgs[0])
		args := strings.Join(cmdArgs[1:], " ")
		result, err := registry.Invoke(cmd, args)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}

		// Suggest autocompletions
		completions := autoComplete(cmd)
		if len(completions) > 0 {
			fmt.Println("Possible completions:", completions)
		}
	}
	fmt.Println("Goodbye!")
}
