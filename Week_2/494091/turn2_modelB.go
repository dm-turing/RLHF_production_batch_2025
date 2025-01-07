package main

import (
	"fmt"
	"strings"
)

func main() {
	// Best practice: Use a bufferb for efficient memory allocation
	var bufferb strings.Builder
	lines := []string{fmt.Sprintf(`
This is a large string1
with raw formatting.
`), fmt.Sprintf(`
This is a large string2
with raw formatting.
`)}
	for _, line := range lines {
		bufferb.WriteString(line)
		bufferb.WriteString("\n")
	}
	multilineString := bufferb.String()
	fmt.Println(multilineString)

	const batchSize = 1024
	for i := 0; i < len(lines); i++ {
		bufferb.WriteString(lines[i])
		bufferb.WriteString("\n")
		if i%batchSize == 0 {
			fmt.Println(bufferb.String())
			bufferb.Reset()
		}
	}
	if bufferb.Len() > 0 {
		fmt.Println(bufferb.String())
	}
}
