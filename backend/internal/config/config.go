package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	JWTSecret          string
	AccessTokenTTLMin  int
	RefreshTokenTTLDays int
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	FrontendOrigin     string
	CookieSecure       bool
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPass           string
	FromEmail          string
	RajaOngkirAPIKey   string
	RajaOngkirBase     string
}

func Load() *Config {
	_ = godotenv.Load()

	smtpPort := 587
	if p, err := strconv.Atoi(getEnv("SMTP_PORT", "587")); err == nil {
		smtpPort = p
	}

	accessTTL := 15
	if t, err := strconv.Atoi(getEnv("ACCESS_TOKEN_TTL_MIN", "15")); err == nil {
		accessTTL = t
	}

	refreshTTL := 7
	if t, err := strconv.Atoi(getEnv("REFRESH_TOKEN_TTL_DAYS", "7")); err == nil {
		refreshTTL = t
	}

	redisDB := 0
	if db, err := strconv.Atoi(getEnv("REDIS_DB", "0")); err == nil {
		redisDB = db
	}

	cfg := &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-change-me"),
		AccessTokenTTLMin:  accessTTL,
		RefreshTokenTTLDays: refreshTTL,
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "ecommerce_auth"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		FrontendOrigin:     getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		CookieSecure:       getEnv("COOKIE_SECURE", "false") == "true",
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            redisDB,
		SMTPHost:           getEnv("SMTP_HOST", ""),
		SMTPPort:           smtpPort,
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPass:           getEnv("SMTP_PASS", ""),
		FromEmail:          getEnv("FROM_EMAIL", "no-reply@ecommerce.local"),
		RajaOngkirAPIKey:   getEnv("RAJAONGKIR_API_KEY", ""),
		RajaOngkirBase:     getEnv("RAJAONGKIR_BASE_URL", "https://api.rajaongkir.com/starter"),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must not be empty")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}
