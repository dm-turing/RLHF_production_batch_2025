package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Strings Library REPL. Type 'help' for a list of commands.")

	for {
		fmt.Print("repl> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			break // Exit if there's an error reading from input
		}

		if string(line) == "quit" || string(line) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Parse and execute the command
		result := executeCommand(string(line))

		if len(result) > 0 {
			fmt.Println(result)
		}
	}
}

func executeCommand(command string) string {
	// Strip leading and trailing whitespace and split by space to get command parts
	parts := strings.Fields(command)

	// Help command
	if len(parts) >= 1 && parts[0] == "help" {
		return printHelp()
	}

	// Handle missing or invalid commands
	if len(parts) == 0 || strings.HasPrefix(command, "//") || strings.HasPrefix(command, "/*") {
		return "" // Do nothing on empty, comment-only lines
	}

	// Match the first command part to determine which string function to use
	switch parts[0] {
	case "ToUpper":
		if len(parts) != 2 {
			return "Usage: ToUpper <string>"
		}
		return strings.ToUpper(parts[1])
	case "ToLower":
		if len(parts) != 2 {
			return "Usage: ToLower <string>"
		}
		return strings.ToLower(parts[1])
	case "Contains":
		if len(parts) != 3 {
			return "Usage: Contains <string> <substring>"
		}
		return fmt.Sprintf("%v contains %v: %v", parts[1], parts[2], strings.Contains(parts[1], parts[2]))
	case "ContainsAny":
		if len(parts) != 3 {
			return "Usage: ContainsAny <string> <substring>"
		}
		return fmt.Sprintf("%v contains any of %v: %v", parts[1], parts[2], strings.ContainsAny(parts[1], parts[2]))
	case "Replace":
		if len(parts) != 4 {
			return "Usage: Replace <string> <old> <new> <count>"
		}
		count, err := strconv.Atoi(parts[3])
		if err != nil {
			return "Usage: Replace <string> <old> <new> <count>"
		}
		return strings.Replace(parts[1], parts[2], parts[3], count)
	case "Join":
		if len(parts) < 3 {
			return "Usage: Join <separator> <string1> <string2> ..."
		}
		return strings.Join(parts[1:], parts[0])
	case "Split":
		if len(parts) != 3 {
			return "Usage: Split <separator> <string>"
		}
		return fmt.Sprintf("%v", strings.Split(parts[2], parts[1]))
	default:
		return fmt.Sprintf("Unknown command: %v. Type 'help' for more information.", parts[0])
	}
}

func printHelp() string {
	return `
Supported Commands:
- ToUpper <string>: Converts all characters in the string to uppercase.
- ToLower <string>: Converts all characters in the string to lowercase.
- Contains <string> <substring>: Returns true if the substring is in the string, false otherwise.
- ContainsAny <string> <substring>: Returns true if any of the characters in the substring is in the string, false otherwise.
- Replace <string> <old> <new> <count>: Replaces all occurrences of the old substring with the new substring in the string, up to 'count' times.
- Join <separator> <string1> <string2> ...: Joins the given strings into a single string separated by the specified separator.
- Split <separator> <string>: Splits the string into slices where the string is separated by the given separator.
- help: Shows this help message.
- quit: Exits the REPL.
`
}
