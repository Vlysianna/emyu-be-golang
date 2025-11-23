package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string
	Port     int
	Env      string
	JWTKey   string
	AppName  string
}

var AppConfig Config

func LoadConfig() error {
	godotenv.Load()

	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	AppConfig = Config{
		DBHost:   getEnv("DB_HOST", "127.0.0.1"),
		DBPort:   dbPort,
		DBUser:   getEnv("DB_USER", "root"),
		DBPass:   getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "emyu"),
		Port:     port,
		Env:      getEnv("SERVER_ENV", "development"),
		JWTKey:   getEnv("JWT_SECRET", "secret"),
		AppName:  getEnv("APP_NAME", "Emyu E-Commerce API"),
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func (c Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
