package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	Env                     string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBSslMode               string
	DBTimeZone              string
	DBAutoMigrate           bool
	SupabaseURL             string
	SupabaseServiceRoleKey  string
	SupabaseStorageBucket   string
	SupabaseChatBucket      string
	SupabaseChatMaxFileSize int64
	FirebaseDatabaseURL     string
	FirebaseDatabaseSecret  string
	JWTSecret               string
	JWTRefreshSecret        string
	JWTAccessExpiration     int
	JWTRefreshExpiration    int
	MaxActiveRefreshTokens  int
	CORSAllowedOrigins      string
	FrontendURL             string
	SMTPHost                string
	SMTPPort                string
	SMTPUsername            string
	SMTPPassword            string
	SMTPFrom                string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system environment variables.")
	}

	return &Config{
		Port:                    getEnv("PORT", "8080"),
		Env:                     getEnv("ENV", "development"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBName:                  getEnv("DB_NAME", "nalakarsa"),
		DBSslMode:               getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:              getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		DBAutoMigrate:           getEnvAsBool("DB_AUTO_MIGRATE", true),
		SupabaseURL:             getEnv("SUPABASE_URL", ""),
		SupabaseServiceRoleKey:  getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		SupabaseStorageBucket:   getEnv("SUPABASE_STORAGE_BUCKET", ""),
		SupabaseChatBucket:      getEnv("SUPABASE_CHAT_BUCKET", "chat-attachments"),
		SupabaseChatMaxFileSize: int64(getEnvAsInt("SUPABASE_CHAT_MAX_FILE_SIZE", 10485760)),
		FirebaseDatabaseURL:     getEnv("FIREBASE_DATABASE_URL", ""),
		FirebaseDatabaseSecret:  getEnv("FIREBASE_DATABASE_SECRET", ""),
		JWTAccessExpiration:     getEnvAsInt("JWT_ACCESS_EXPIRATION", 900),
		JWTRefreshExpiration:    getEnvAsInt("JWT_REFRESH_EXPIRATION", 604800),
		MaxActiveRefreshTokens:  getEnvAsInt("MAX_ACTIVE_REFRESH_TOKENS", 5),
		CORSAllowedOrigins:      getEnv("CORS_ALLOWED_ORIGINS", "*"),
		FrontendURL:             getEnv("FRONTEND_URL", "http://localhost:5173"),
		SMTPHost:                getEnv("SMTP_HOST", ""),
		SMTPPort:                getEnv("SMTP_PORT", "587"),
		SMTPUsername:            getEnv("SMTP_USERNAME", ""),
		SMTPPassword:            getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                getEnv("SMTP_FROM", ""),
		DBPassword:       getEnvRequired("DB_PASSWORD"),
		JWTSecret:        getEnvRequired("JWT_SECRET"),
		JWTRefreshSecret: getEnvRequired("JWT_REFRESH_SECRET"),
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set.", key)
	}
	return value
}
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
