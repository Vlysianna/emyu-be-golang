package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/emyu/ecommer-be/config"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("mysql", config.AppConfig.GetDSN())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return err
	}

	log.Println("✓ Database connected successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
