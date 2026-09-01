package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
		DatabaseURL:           normalizeDatabaseURL(getEnv("DATABASE_URL", "")),
		MongoDatabase:         getEnv("MONGO_DATABASE", "travel_diary"),
		AWSRegion:             getEnv("AWS_REGION", ""),
		AWSAccessKeyID:        getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:    getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSBucket:             getEnv("AWS_S3_BUCKET", ""),
		JWTSecret:             getEnv("JWT_SECRET", "change-me"),
		CORSAllowedOrigins:    sanitizeOrigins(getEnv("CORS_ALLOWED_ORIGINS", "")),
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleAuthStateSecret: getEnv("GOOGLE_AUTH_STATE_SECRET", "state-secret"),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL", ""),
		FrontendRedirectURL:   getEnv("FRONTEND_REDIRECT_URL", ""),
	}
}

func normalizeDatabaseURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, `"'`)
	raw = strings.TrimPrefix(raw, "export DATABASE_URL=")
	raw = strings.TrimPrefix(raw, "DATABASE_URL=")
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	return raw
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
		value = sanitizeEnvValue(strings.TrimSpace(value))
		if key == "" || value == "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

var markdownLinkPattern = regexp.MustCompile(`^\[([^\]]+)\]\([^)]+\)$`)

func sanitizeEnvValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	if match := markdownLinkPattern.FindStringSubmatch(value); len(match) == 2 {
		return match[1]
	}
	return value
}

func sanitizeOrigins(raw string) string {
	parts := strings.Split(raw, ",")
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		part = sanitizeEnvValue(part)
		if part == "" {
			continue
		}
		if u, err := url.Parse(part); err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
			clean = append(clean, part)
		}
	}
	return strings.Join(clean, ",")
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

func ValidateDatabaseURL(raw string) error {
	if raw == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if parsed.Scheme != "mongodb" && parsed.Scheme != "mongodb+srv" {
		return fmt.Errorf("DATABASE_URL must start with mongodb:// or mongodb+srv://")
	}
	return nil
}
