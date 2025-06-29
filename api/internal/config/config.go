package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	JWTKey         string
	JWTRefreshKey  string
	JWTAccessTime  string
	JWTRefreshTime string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", ""),
		DBPass:         getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", ""),
		JWTKey:         getEnv("JWT_SECRET", ""),
		JWTRefreshKey:  getEnv("JWT_SECRET_REFRESH", ""),
		JWTAccessTime:  getEnv("JWT_ACCESS_EXPIRES", ""),
		JWTRefreshTime: getEnv("JWT_REFRESH_EXPIRES", ""),
	}
}

func getEnv(key, fallback string) string {
	log.Println(key)
	if val, ok := os.LookupEnv(key); ok {
		log.Println(val)
		return val
	}
	return fallback
}
