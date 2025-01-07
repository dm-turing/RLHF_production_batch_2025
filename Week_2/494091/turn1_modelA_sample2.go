package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {
	name := "Alice"
	age := 30
	city := "New York"

	tmpl := template.Must(template.New("").Parse(`
    Name: {{.Name}}
    Age: {{.Age}}
    City: {{.City}}
    `))

	data := map[string]interface{}{
		"Name": name,
		"Age":  age,
		"City": city,
	}

	var multilineString bytes.Buffer
	err := tmpl.Execute(&multilineString, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println(multilineString.String())
}
