package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	AccessTokenTTLMin  int
	RefreshTokenTTLDays int
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	FrontendOrigin     string
	CookieSecure       bool
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults()

	_ = viper.ReadInConfig()

	cfg := &Config{
		AppPort:             viper.GetString("APP_PORT"),
		DBHost:              viper.GetString("DB_HOST"),
		DBPort:              viper.GetString("DB_PORT"),
		DBUser:              viper.GetString("DB_USER"),
		DBPassword:          viper.GetString("DB_PASSWORD"),
		DBName:              viper.GetString("DB_NAME"),
		DBSSLMode:           viper.GetString("DB_SSLMODE"),
		JWTSecret:           viper.GetString("JWT_SECRET"),
		AccessTokenTTLMin:   viper.GetInt("ACCESS_TOKEN_TTL_MIN"),
		RefreshTokenTTLDays: viper.GetInt("REFRESH_TOKEN_TTL_DAYS"),
		RedisAddr:           viper.GetString("REDIS_ADDR"),
		RedisPassword:       viper.GetString("REDIS_PASSWORD"),
		RedisDB:             viper.GetInt("REDIS_DB"),
		FrontendOrigin:      viper.GetString("FRONTEND_ORIGIN"),
		CookieSecure:        strings.EqualFold(viper.GetString("COOKIE_SECURE"), "true"),
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("incomplete database configuration")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("jwt secret must not be empty")
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "ecommerce_auth")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "super-secret-change-me")
	viper.SetDefault("ACCESS_TOKEN_TTL_MIN", "15")
	viper.SetDefault("REFRESH_TOKEN_TTL_DAYS", "7")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", "0")
	viper.SetDefault("FRONTEND_ORIGIN", "http://localhost:5173")
	viper.SetDefault("COOKIE_SECURE", "false")
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

func (c *Config) RedisDBAsString() string {
	return strconv.Itoa(c.RedisDB)
}
