package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"visualon.com/go-server-monitor/config"
)

// DB Global DB connection
var DB *sql.DB

func init() {
	var err error

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.CONFIG.Database.DBUser, config.CONFIG.Database.DBPassword, config.CONFIG.Database.DBHost, config.CONFIG.Database.DBPort, config.CONFIG.Database.DBName)
	db, err := sql.Open("mysql", connectionString)
	DB = db
	if err != nil {
		panic(err)
	}
}
