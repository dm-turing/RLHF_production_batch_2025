package main

import (
	"fmt"
	"reflect"
)

// Define your model structs
type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

type Product struct {
	ID     int     `db:"id"`
	Name   string  `db:"name"`
	Price  float64 `db:"price"`
	UserID int     `db:"user_id"`
}

// ORM struct using maps
type ORM struct {
	models map[string]map[int]interface{}
}

// NewORM creates a new ORM instance
func NewORM() *ORM {
	return &ORM{
		models: make(map[string]map[int]interface{}),
	}
}

// Insert inserts a new model into the ORM
func (orm *ORM) Insert(model interface{}) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Invalid model type: %s", t.Kind())
	}

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Model must be a pointer: %s", t.Kind())
	}

	v = v.Elem()

	modelName := t.Name()
	if _, ok := orm.models[modelName]; !ok {
		orm.models[modelName] = make(map[int]interface{})
	}

	// Get the primary key value using the "db" tag
	// primaryKeyTag := t.Field(0).Tag.Get("db")
	primaryKeyValue := v.Field(0).Int()

	// Store the model in the ORM
	orm.models[modelName][int(primaryKeyValue)] = model

	return nil
}

// FindByID finds a model by its primary key
func (orm *ORM) FindByID(modelName string, id int) interface{} {
	if _, ok := orm.models[modelName]; !ok {
		return nil
	}

	return orm.models[modelName][id]
}

// Delete deletes a model by its primary key
func (orm *ORM) Delete(modelName string, id int) {
	if _, ok := orm.models[modelName]; ok {
		delete(orm.models[modelName], id)
	}
}

func main() {
	orm := NewORM()

	// Create new users and products
	user1 := &User{ID: 1, Name: "Alice", Age: 25}
	user2 := &User{ID: 2, Name: "Bob", Age: 30}
	product1 := &Product{ID: 1, Name: "Phone", Price: 699.99, UserID: 1}
	product2 := &Product{ID: 2, Name: "Laptop", Price: 999.99, UserID: 2}

	// Insert models into the ORM
	err := orm.Insert(user1)
	if err != nil {
		panic(err)
	}
	err = orm.Insert(user2)
	if err != nil {
		panic(err)
	}
	err = orm.Insert(product1)
	if err != nil {
		panic(err)
	}
	err = orm.Insert(product2)
	if err != nil {
		panic(err)
	}

	// Find a user by ID
	foundUser := orm.FindByID("User", 1)
	fmt.Println("Found User:", foundUser.(*User).Name)

	// Delete a product by ID
	orm.Delete("Product", 2)
}
