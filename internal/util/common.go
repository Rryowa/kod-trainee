package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kod/internal/models"
	"kod/internal/models/config"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ParseUserForm(r *http.Request) (*models.User, error) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	user := &models.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if len(user.Username) == 0 {
		return nil, ErrEmptyUsername
	} else if len(user.Password) == 0 {
		return nil, ErrEmptyPassword
	}

	return user, nil
}

func GetUserIdFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return -1, errors.New("you are unauthorized")
	}
	return userID, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		HandleHttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func WriteJSONToken(w http.ResponseWriter, token string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
}

func NewHttpConfig() *config.HttpConfig {
	return &config.HttpConfig{
		Host:      os.Getenv("HTTP_HOST"),
		Port:      os.Getenv("HTTP_PORT"),
		JwtSecret: os.Getenv("JWT_SECRET"),
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