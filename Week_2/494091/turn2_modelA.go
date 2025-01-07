package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	multilineString := fmt.Sprintf(`
This is a multiline string
with raw formatting.
`)
	fmt.Println(multilineString)
	largeData := []string{fmt.Sprintf(`
This is a large string1
with raw formatting.
`), fmt.Sprintf(`
This is a large string2
with raw formatting.
`)}

	var builder strings.Builder
	for _, line := range largeData {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	multilineString = builder.String()
	fmt.Println(multilineString)

	var buffer bytes.Buffer
	for _, line := range largeData {
		fmt.Fprintln(&buffer, line)
	}
	multilineString = buffer.String()
	fmt.Println(multilineString)

	var buf bytes.Buffer
	prefix := "Largedata:"
	buf.WriteString(fmt.Sprintf("Prefix: %s", prefix))
	for _, line := range largeData {
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	multilineString = buf.String()
	fmt.Println(multilineString)
}
