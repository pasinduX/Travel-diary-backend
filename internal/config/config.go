package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppName               string
	Environment           string
	Port                  string
	DatabaseURL           string
	MongoDatabase         string
	AWSRegion             string
	AWSAccessKeyID        string
	AWSSecretAccessKey    string
	AWSBucket             string
	JWTSecret             string
	CORSAllowedOrigins    string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleAuthStateSecret string
	GoogleRedirectURL     string
	FrontendRedirectURL   string
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		AppName:               getEnv("APP_NAME", "travel-diary-backend"),
		Environment:           getEnv("APP_ENV", "development"),
		Port:                  getEnv("PORT", "8001"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		MongoDatabase:         getEnv("MONGO_DATABASE", "travel_diary"),
		AWSRegion:             getEnv("AWS_REGION", ""),
		AWSAccessKeyID:        getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:    getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSBucket:             getEnv("AWS_S3_BUCKET", ""),
		JWTSecret:             getEnv("JWT_SECRET", "change-me"),
		CORSAllowedOrigins:    getEnv("CORS_ALLOWED_ORIGINS", ""),
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleAuthStateSecret: getEnv("GOOGLE_AUTH_STATE_SECRET", "state-secret"),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL", ""),
		FrontendRedirectURL:   getEnv("FRONTEND_REDIRECT_URL", ""),
	}
}

func loadDotEnv(path string) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" || value == "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func (c Config) ServerAddress() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
