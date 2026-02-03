package utils

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	DatabaseIp   = GetEnv("DB_MIDDLEWARE_IP", "10.254.254.110")
	DatabasePort = GetEnv("DB_MIDDLEWARE_PORT", "3306")
	DatabaseUser = GetEnv("DB_MIDDLEWARE_USER", "jpa")
	DatabasePass = GetEnv("DB_MIDDLEWARE_PASS", "SRsUgEwKi7cL84VnuxP8pGyfzzvtkJ1LRqEIfQBarl+gW6m6Y22rXed4Ras=")
	DatabaseName = GetEnv("DB_MIDDLEWARE_NAME", "DB_MID_REKON")
	AesEncKey    = GetEnv("AES_ENC_KEY", "VAUNkTWRJlCOnXKPNhpU1w==")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		// fmt.Printf("[ENV] %s=%s \n", key, value)
		return value
	}
	// fmt.Printf("[FALLBACK] %s=%s \n", key, fallback)
	return fallback
}
