package main

import (
	"fmt"
)

func main() {
	multilineString := fmt.Sprintf(`  
This is a multiline string  
that can be used for logging purposes  
or any other reason where a multi-line  
string is required.  
`)

	fmt.Println(multilineString)
}
