package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func DBConn() (db *sql.DB) {
	user := os.Getenv("GW_DB_USER")
	pass := os.Getenv("GW_DB_PASS")
	host := os.Getenv("GW_DB_HOST")
	port := os.Getenv("GW_DB_PORT")
	dbname := os.Getenv("GW_DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}
