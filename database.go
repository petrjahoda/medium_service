package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const config = "user=postgres password=Medium dbname=postgres host=medium_database port=5432 sslmode=disable"
const mediumConfig = "user=postgres password=Medium dbname=medium host=medium_database port=5432 sslmode=disable"

type ButtonRecords struct {
	gorm.Model
	Name string
	Time time.Time
}

func checkDatabase() {
	time.Sleep(2 * time.Second)
	postgresDb, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	postgresSqlDB, _ := postgresDb.DB()
	defer postgresSqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening postgres database: "+err.Error())
		return
	}
	logInfo("MAIN", "Postgres  database available, checking for medium database")
	mediumDb, err := gorm.Open(postgres.Open(mediumConfig), &gorm.Config{})
	mediumSqlDB, _ := mediumDb.DB()
	defer mediumSqlDB.Close()
	if err != nil {
		logError("MAIN", "Medium database does not exist, creating")
		postgresDb = postgresDb.Exec("CREATE DATABASE medium;")
		if postgresDb.Error != nil {
			logError("MAIN", "Cannot create medium database: "+err.Error())
			return
		}
		logInfo("MAIN", "Medium database created")
	}
	logInfo("MAIN", "Medium  database available")

	if !mediumDb.Migrator().HasTable(&ButtonRecords{}) {
		logInfo("MAIN", "ButtonRecords table does not exist, creating")
		err := mediumDb.Migrator().CreateTable(&ButtonRecords{})
		if err != nil {
			logError("MAIN", "Cannot create ButtonRecords table")
			return
		}
	} else {
		logInfo("MAIN", "ButtonRecords table exists, updating")
		err := mediumDb.Migrator().AutoMigrate(&ButtonRecords{})
		if err != nil {
			logError("MAIN", "Cannot update ButtonRecords table")
			return
		}
		logInfo("MAIN", "ButtonRecords table updated")
	}
}
