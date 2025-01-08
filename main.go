package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Query represents a parsed query part.
type Query struct {
	Field    string
	Value    string
	Operator string
}

func parseQuery(input string) (queries []Query) {
	termRegex := regexp.MustCompile(`\b([a-zA-Z0-9_]+):([a-zA-Z0-9_]+)\b`)
	clauseRegex := regexp.MustCompile(`(AND|OR)`)

	// Split the input string by operators
	clauses := clauseRegex.Split(input, -1)
	for _, clause := range clauses {
		terms := termRegex.FindAllStringSubmatch(clause, -1)
		if len(terms) > 0 {
			for _, term := range terms {
				queries = append(queries, Query{
					Field:    term[1],
					Value:    term[2],
					Operator: strings.ToLower(clauseRegex.FindStringSubmatch(input)[1]),
				})
			}
		} else if clause != "" {
			// Handle remaining operators (if at the end or nothing in between)
			queries[len(queries)-1].Operator = strings.ToLower(clauseRegex.FindStringSubmatch(input)[1])
		}
	}

	return
}

func reconstructURL(queries []Query) string {
	parts := make([]string, 0, len(queries))
	var lastOp string

	for _, query := range queries {
		parts = append(parts, fmt.Sprintf("%s=%s", query.Field, query.Value))
		lastOp = query.Operator
	}

	if lastOp == "and" {
		parts = append(parts, "&")
	} else if lastOp == "or" {
		parts = append(parts, "|")
	}

	return strings.Join(parts, "&")
}

func main() {
	inputQuery := "name:John AND age:30 OR location:New York"

	queries := parseQuery(inputQuery)

	fmt.Println("Parsed Queries:")
	for _, query := range queries {
		fmt.Printf("%#v\n", query)
	}

	reconstructedURL := reconstructURL(queries)
	fmt.Printf("\nReconstructed URL: %s\n", reconstructedURL)
}
