package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	GinMode              string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	JWTSecret            string
	JWTExpirationHours   int
	CORSOrigins          string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetimeMin int
	RateLimitRPS         float64
	RateLimitBurst       float64
	RequestTimeoutSec    int
}

var App Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env não encontrado — usando variáveis de ambiente do sistema")
	}

	hours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		hours = 24
	}

	dbMaxOpen, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		dbMaxOpen = 25
	}

	dbMaxIdle, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		dbMaxIdle = 5
	}

	dbConnMaxLifetime, err := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME_MIN", "30"))
	if err != nil {
		dbConnMaxLifetime = 30
	}

	rateLimitRPS, err := strconv.ParseFloat(getEnv("RATE_LIMIT_RPS", "20"), 64)
	if err != nil {
		rateLimitRPS = 20
	}

	rateLimitBurst, err := strconv.ParseFloat(getEnv("RATE_LIMIT_BURST", "50"), 64)
	if err != nil {
		rateLimitBurst = 50
	}

	requestTimeoutSec, err := strconv.Atoi(getEnv("REQUEST_TIMEOUT_SEC", "30"))
	if err != nil {
		requestTimeoutSec = 30
	}

	App = Config{
		Port:                 getEnv("PORT", "8080"),
		GinMode:              getEnv("GIN_MODE", "debug"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "endurance"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		JWTSecret:            getEnv("JWT_SECRET", "secret_dev_only"),
		JWTExpirationHours:   hours,
		CORSOrigins:          getEnv("CORS_ORIGINS", "http://localhost:5173"),
		DBMaxOpenConns:       dbMaxOpen,
		DBMaxIdleConns:       dbMaxIdle,
		DBConnMaxLifetimeMin: dbConnMaxLifetime,
		RateLimitRPS:         rateLimitRPS,
		RateLimitBurst:       rateLimitBurst,
		RequestTimeoutSec:    requestTimeoutSec,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
