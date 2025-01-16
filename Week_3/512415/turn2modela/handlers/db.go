package handlers

import (
	"Week_3/512415/turn2modela/models"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open("sqlite3", "convert_app.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	db.AutoMigrate(&models.Conversion{})
}

func StoreConversion(conversion *models.Conversion) {
	db.Create(conversion)
}
