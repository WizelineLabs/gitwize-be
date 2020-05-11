package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitwize-be/src/configuration"
)

func DBConn() (db *sql.DB) {
	config := configuration.ReadConfiguration()
	user := config.Database.GwDbUser
	pass := config.Database.GwDbPassword
	host := config.Database.GwDbHost
	port := config.Database.GwDbPort
	dbname := config.Database.GwDbName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}
