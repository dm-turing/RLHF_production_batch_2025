package main

import (
	"fmt"
	"log"

	"github.com/misterclunk/peggo"
)

type Query struct {
	Words []string
	Op    string
}

var p *peggo.Parser

func main() {
	query := "name:John | age:30 & (location:New York | location:Los Angeles)"

	// Compile the parser
	p = peggo.NewParser("query.peg")

	// Parse the query
	result, err := p.Parse(query)
	if err != nil {
		log.Fatalf("Error parsing query: %v", err)
	}

	// Output the parsed query as an AST
	fmt.Println("Parsed Query:")
	parseAndOutput(result)

	// Reconstruct the URL from the AST
	fmt.Println("\nReconstructed URL:")
	reconstructURL(result)
}

func parseAndOutput(result interface{}) {
	q, ok := result.(*Query)
	if !ok {
		fmt.Println("Invalid query structure")
		return
	}

	fmt.Printf("Query: %#v\n", q)

	for _, word := range q.Words {
		fmt.Printf("Word: %s\n", word)
	}

	if q.Op != "" {
		fmt.Printf("Operator: %s\n", q.Op)
		for _, sub := range q.Words {
			parseAndOutput(sub)
		}
	}
}

func reconstructURL(result interface{}) string {
	q, ok := result.(*Query)
	if !ok {
		return result.(string)
	}

	var urlComponents []string

	for _, word := range q.Words {
		urlComponents = append(urlComponents, reconstructURL(word))
	}

	if q.Op != "" {
		return fmt.Sprintf("(%s %s %s)", urlComponents[0], q.Op, urlComponents[1])
	}

	return fmt.Sprintf("%s", q.Words[0])
}
