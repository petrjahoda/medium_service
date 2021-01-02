package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const config = "user=postgres password=Medium dbname=postgres host=medium_database port=5432 sslmode=disable"

func checkDatabase() {
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	} else {
		logInfo("MAIN", "Database available!")
	}
}
