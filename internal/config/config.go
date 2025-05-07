package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)


var (
	DB_URL             string
	PORT               string
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	REDIS_URL          string
	PRODUCTION         bool
)

var once     sync.Once

func LoadConfig() {

	once.Do(func() {
		loadEnv()

		DB_URL=    			getEnv("DB_URL", "postgresql://postgres:123@db:5432/userm")
		JWT_ACCESS_SECRET= 	getEnv("JWT_ACCESS_SECRET", "mysecret")
		JWT_REFRESH_SECRET= getEnv("JWT_REFRESH_SECRET", "mysecretrefresh")
		PORT= 				getEnv("PORT", "8080")
		REDIS_URL= 			getEnv("REDIS_URL", "redis:6379")
		PRODUCTION= 		getEnv("PRODUCTION", "false") == "true"
	})

}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback

}