package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                       string
	Env                        string
	DBHost                     string
	DBPort                     string
	DBUser                     string
	DBPassword                 string
	DBName                     string
	DBSslMode                  string
	DBTimeZone                 string
	FirebaseProjectID          string
	FirebaseCredentialJSONPath string
	FirebaseStorageBucket      string
	JWTSecret                  string
	JWTRefreshSecret           string
	JWTAccessExpiration        int
	JWTRefreshExpiration       int
	CORSAllowedOrigins         string
	FrontendURL                string
}

func LoadConfig() *Config {
	// Memuat file .env (untuk dev lokal)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system environment variables.")
	}

	return &Config{
		Port:                       getEnv("PORT", "8080"),
		Env:                        getEnv("ENV", "development"),
		DBHost:                     getEnv("DB_HOST", "localhost"),
		DBPort:                     getEnv("DB_PORT", "5432"),
		DBUser:                     getEnv("DB_USER", "postgres"),
		DBName:                     getEnv("DB_NAME", "nalakarsa"),
		DBSslMode:                  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:                 getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		FirebaseProjectID:          getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseCredentialJSONPath: getEnv("FIREBASE_CREDENTIAL_JSON_PATH", ""),
		FirebaseStorageBucket:      getEnv("FIREBASE_STORAGE_BUCKET", ""),
		JWTAccessExpiration:        getEnvAsInt("JWT_ACCESS_EXPIRATION", 900),
		JWTRefreshExpiration:       getEnvAsInt("JWT_REFRESH_EXPIRATION", 604800),
		CORSAllowedOrigins:         getEnv("CORS_ALLOWED_ORIGINS", "*"),
		FrontendURL:                getEnv("FRONTEND_URL", "http://localhost:5173"),

		// variabel rahasia (wajib diisi dari .env)
		DBPassword:       getEnvRequired("DB_PASSWORD"),
		JWTSecret:        getEnvRequired("JWT_SECRET"),
		JWTRefreshSecret: getEnvRequired("JWT_REFRESH_SECRET"),
	}
}

// fungsi untuk variabel opsional
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// fungsi untuk variabel wajib
func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set.", key)
	}
	return value
}

// Fungsi untuk Konversi Integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
