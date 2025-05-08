package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)


var (
	DBUrl             	string
	Port               	string
	JWTAccessSecret	  	string
	JWTRefreshSecret 	string
	RedisURL          	string
	Production         	bool
)

var once     sync.Once

func LoadConfig() {

	once.Do(func() {
		loadEnv()

		DBUrl=    			getEnv("DB_URL", "postgresql://postgres:123@db:5432/userm")
		JWTAccessSecret= 	getEnv("JWT_ACCESS_SECRET", "mysecret")
		JWTRefreshSecret= 	getEnv("JWT_REFRESH_SECRET", "mysecretrefresh")
		Port= 				getEnv("PORT", "8080")
		RedisURL= 			getEnv("REDIS_URL", "redis:6379")
		Production= 		getEnv("PRODUCTION", "false") == "true"
	})

}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Info("No .env file found")
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback

}