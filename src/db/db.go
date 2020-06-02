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

var gormDB *gorm.DB

func dbConn() (db *gorm.DB) {
	user := configuration.CurConfiguration.Database.GwDbUser
	pass := configuration.CurConfiguration.Database.GwDbPassword
	host := configuration.CurConfiguration.Database.GwDbHost
	port := configuration.CurConfiguration.Database.GwDbPort
	dbname := configuration.CurConfiguration.Database.GwDbName

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	return
}

func Initialize() {
	gormDB = dbConn()

	// Migrate the schema
	gormDB.AutoMigrate(&Repository{}, &Metric{})
}

func SqlDBConn() (db *sql.DB) {
	dbConn := os.Getenv("DB_CONN_STRING")

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	return
}
