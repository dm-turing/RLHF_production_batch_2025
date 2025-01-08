package main

import (
	"fmt"
	"strings"
)

type Query struct {
	Text       string   `json:"text,omitempty"`
	Categories []string `json:"categories,omitempty"`
	PriceRange struct {
		Min float64
		Max float64
	} `json:"price_range,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (q *Query) ToURL() string {
	var parts []string

	// Add text query
	if q.Text != "" {
		parts = append(parts, fmt.Sprintf("q=%s", q.Text))
	}

	// Add categories
	if len(q.Categories) > 0 {
		catStr := strings.Join(q.Categories, ",")
		parts = append(parts, fmt.Sprintf("categories=%s", catStr))
	}

	// Add price range
	if q.PriceRange.Min > 0 && q.PriceRange.Max > 0 {
		parts = append(parts, fmt.Sprintf("price_range=%f-%f", q.PriceRange.Min, q.PriceRange.Max))
	}

	// Add attributes
	for attr, val := range q.Attributes {
		parts = append(parts, fmt.Sprintf("attributes[%s]=%s", attr, val))
	}

	return "/search?" + strings.Join(parts, "&")
}

func main() {
	q := Query{
		Text:       "Laptop",
		Categories: []string{"Electronics", "Computers"},
		PriceRange: struct {
			Min float64
			Max float64
		}{float64(100), float64(1000)},
		Attributes: map[string]string{
			"Brand":       "Dell",
			"Screen_Size": "15.6",
		},
	}

	fmt.Println(q.ToURL())
	// Output: /search?q=Laptop&categories=Electronics,Computers&price_range=100.00-1000.00&attributes[Brand]=Dell&attributes[Screen_Size]=15.6
}
