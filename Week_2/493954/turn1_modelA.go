package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

type Product struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Categories  []string       `json:"categories"`
	Price       float64        `json:"price"`
	Rating      float64        `json:"rating"`
	Stock       int            `json:"stock"`
	Attributes  map[string]any `json:"attributes"` // Allow for dynamic attributes
}

func CreateIndex(client *elastic.Client) error {
	// Create a new index
	ctx := context.Background()
	_, err := client.Index().Index("ecommerce").BodyString(`
	{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"title": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"description": {
					"type": "text"
				},
				"categories": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"price": {
					"type": "float"
				},
				"rating": {
					"type": "float"
				},
				"stock": {
					"type": "long"
				},
				"attributes": {
					"type": "object",
					"dynamic": "true"
				}
			}
		}
	}
	`).Do(ctx)
	return err
}

func IndexProduct(client *elastic.Client, product *Product) error {
	ctx := context.Background()
	_, err := client.Index().Index("ecommerce").Id(product.ID).BodyJson(product).Do(ctx)
	return err
}

func SearchProducts(client *elastic.Client, query string) ([]Product, error) {
	ctx := context.Background()
	var products []Product

	esSearch := client.Search().Index("ecommerce")

	// Simple example of a search query parsing; in practice, you'll want a robust query builder
	esSearch.Query(elastic.NewQueryStringQuery(query))

	searchResp, err := esSearch.Do(ctx)
	if err != nil {
		log.Printf("Error searching: %v", err)
		return products, err
	}

	for _, hit := range searchResp.Hits.Hits {
		var product Product
		err = json.Unmarshal(*hit.Source, &product)
		if err != nil {
			log.Printf("Error unmarshaling product: %v", err)
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

func main() {
	c, _ := elastic.NewClient()
	err := CreateIndex(c)
	if err != nil {
		fmt.Println(err)
	}
}
