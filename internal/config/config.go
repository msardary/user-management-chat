package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)


type ConfigStruct struct {
	DB_URL 				string
	PORT  				string
	JWT_ACCESS_SECRET 	string
	JWT_REFRESH_SECRET 	string
	REDIS_URL 			string
}

var Load ConfigStruct

func LoadConfig() {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Load = ConfigStruct{
		DB_URL:    			getEnv("DB_URL", "postgresql://postgres:123@db:5432/userm"),
		JWT_ACCESS_SECRET: 	getEnv("JWT_ACCESS_SECRET", "mysecret"),
		JWT_REFRESH_SECRET: getEnv("JWT_REFRESH_SECRET", "mysecretrefresh"),
		PORT: 				getEnv("PORT", "8080"),
		REDIS_URL: 			getEnv("REDIS_URL", "redis:6379"),
	}

}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback

}