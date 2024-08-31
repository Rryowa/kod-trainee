package util

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kod/internal/models/config"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
}

func NewHttpConfig() *config.HttpConfig {
	return &config.HttpConfig{
		Host:          os.Getenv("HTTP_HOST"),
		Port:          os.Getenv("HTTP_PORT"),
		TelemetryAddr: os.Getenv("TELEMETRY_ADDR"),
	}
}

func NewSessionConfig() *config.SessionConfig {
	cookieTtl, err := time.ParseDuration(os.Getenv("COOKIE_TTL"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}
	jwtTtl, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &config.SessionConfig{
		CookieTTL:  cookieTtl,
		CookieName: os.Getenv("COOKIE_NAME"),
		JwtTTL:     jwtTtl,
		JwtSecret:  os.Getenv("JWT_SECRET"),
	}
}

func NewDbConfig() *config.DbConfig {
	attempts, err := strconv.Atoi(os.Getenv("ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting ATTEMPTS: %v\n", err)
	}
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &config.DbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Attempts: attempts,
		Timeout:  timeout,
	}
}

func NewZapLogger() *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	sugar := logger.Sugar()
	sugar.Sync()
	return sugar
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}