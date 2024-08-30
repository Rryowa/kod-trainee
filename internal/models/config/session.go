package config

import "time"

type SessionConfig struct {
	CookieTTL  time.Duration `env:"COOKIE_TTL" envDefault:"10m"`
	CookieName string        `env:"COOKIE_NAME" envDefault:"jwt"`
	JwtTTL     time.Duration `env:"JWT_TTL" envDefault:"5m"`
	JwtSecret  string        `env:"JWT_SECRET"`
}