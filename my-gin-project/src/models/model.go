package models

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@tcp(" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrer les modèles
	db.AutoMigrate(&User{}, &Item{})

	// Conserver la DB globale
	DB = db

	return db, nil
}

func (item *Item) UpdatePrice(newPrice float64) {
	item.Price = newPrice
}
