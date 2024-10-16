package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Name        string
	Tags        []string               `gorm:"type:bytes;serializer:json"`
	Spec        map[string]interface{} `gorm:"serializer:json"`
	SpecGob     map[string]interface{} `gorm:"type:bytes;serializer:gob"`
	CreatedTime int64                  `gorm:"serializer:unixtime;type:time"`
}

var DB *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	database, err := gorm.Open(mysql.Open("ehis_dev:devdba_user@tcp(127.0.0.1:3307)/gorm_serializer?charset=utf8&parseTime=true"), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&Product{})
}

func main() {
	connectDatabase()
	dbMigrate()

	createdAt := time.Now()
	data := Product{
		Name: "Apple iPhone 13",
		Tags: []string{"smartphone", "iphone", "apple", "cell phone", "5g", "camera", "retina display"},
		Spec: map[string]interface{}{
			"name":       "Apple iPhone 13",
			"display":    "6.1 inches",
			"resolution": "2532 x 1170 pixels",
			"processor":  "Apple A15 Bionic",
			"ram":        "6GB",
			"storage":    "128GB",
		},
		SpecGob: map[string]interface{}{
			"name":       "Apple iPhone 13",
			"display":    "6.1 inches",
			"resolution": "2532 x 1170 pixels",
			"processor":  "Apple A15 Bionic",
			"ram":        "6GB",
			"storage":    "128GB",
		},
		CreatedTime: createdAt.Unix(),
	}
	DB.Create(&data)

	var result Product
	DB.First(&result, "id = ?", data.ID)

	fmt.Printf("Name: ")
	fmt.Println(result.Name)

	fmt.Printf("Tags: ")
	fmt.Println(result.Tags)

	fmt.Printf("Spec: ")
	fmt.Println(result.Spec)

	fmt.Printf("SpecGob: ")
	fmt.Println(result.SpecGob)

	fmt.Printf("CreatedTime: ")
	fmt.Println(result.CreatedTime)
}
