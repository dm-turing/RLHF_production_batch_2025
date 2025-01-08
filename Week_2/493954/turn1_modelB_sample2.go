
package main

import (
    "fmt"
    "net/url"
)

// SearchQuery represents a structured search query
type SearchQuery struct {
    Text      string                 `json:"text"`
    Attributes map[string]interface{} `json:"attributes"`
    Filters   []Filter               `json:"filters"`
}

// Filter represents a search filter
type Filter struct {
    Key      string   `json:"key"`
    Values   []string `json:"values"`
    Operator string   `json:"operator"`
}

func (q *SearchQuery) ToQueryString() string {
    values := url.Values{}

    if q.Text != "" {
        values.Set("text", q.Text)
    }

    // Handle attributes
    for key, value := range q.Attributes {
        switch v := value.(type) {
        case string:
            values.Set(fmt.Sprintf("attributes[%s]", key), v)
        case map[string]string:
            for subKey, subValue := range v {
                values.Set(fmt.Sprintf("attributes[%s][%s]", key, subKey), subValue)
            }
        default:
            fmt.Printf("Unsupported attribute type: %T\n", v)
        }
    }

    // Handle filters
    for _, filter := range q.Filters {
        for _, value := range filter.Values {
            values.Add(fmt.Sprintf("filters[%s][%s]", filter.Key, filter.Operator), value)
        }
    }

    return values.Encode()
}

func main() {
    // Define a complex search query
    query := SearchQuery{
        Text: "smartphone",
        Attributes: map[string]interface{}{
            "Display.Size":   ">=6.0 inches",
            "Memory":         ">=8GB",
            "Battery.Capacity": ">=4000mAh",