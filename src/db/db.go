package db

import (
	"database/sql"
	"fmt"
	"gitwize-be/src/configuration"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func dbConn() (db *gorm.DB) {
	config := configuration.ReadConfiguration()
	user := config.Database.GwDbUser
	pass := config.Database.GwDbPassword
	host := config.Database.GwDbHost
	port := config.Database.GwDbPort
	dbname := config.Database.GwDbName

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	return
}

func Initialize() *gorm.DB {
	gormDB := dbConn()

	// Migrate the schema
	gormDB.AutoMigrate(&Repository{}, &Metric{})

	return gormDB
}

func SqlDBConn() (db *sql.DB) {
	dbConn := os.Getenv("DB_CONN_STRING")

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	return
}
