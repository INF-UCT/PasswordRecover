package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	App    AppConfig
	LDAP   LDAPConfig
	Redis  RedisConfig
	SMTP   SMTPConfig
	Ramtun RamtunConfig
}

type AppConfig struct {
	Port            string
	BaseURL         string
	FrontendURL     string
	TokenTTLMinutes int
}

type LDAPConfig struct {
	Host         string
	Port         string
	UseTLS       bool
	BindDN       string
	BindPassword string
	BaseDN       string
}

type RedisConfig struct {
	URL      string
	Password string
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type RamtunConfig struct {
	Host   string
	ApiKey string
}

func Load() (*Config, error) {
	var missing []string

	required := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	tokenTTL, _ := strconv.Atoi(getEnvOrDefault("APP_TOKEN_TTL_MINUTES", "10"))
	smtpPort, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
	useTLS, _ := strconv.ParseBool(getEnvOrDefault("LDAP_USE_TLS", "false"))

	cfg := &Config{
		App: AppConfig{
			Port:            getEnvOrDefault("APP_PORT", "8080"),
			BaseURL:         required("APP_BASE_URL"),
			FrontendURL:     required("APP_FRONTEND_URL"),
			TokenTTLMinutes: tokenTTL,
		},
		LDAP: LDAPConfig{
			Host:         required("LDAP_HOST"),
			Port:         getEnvOrDefault("LDAP_PORT", "389"),
			UseTLS:       useTLS,
			BindDN:       required("LDAP_BIND_DN"),
			BindPassword: required("LDAP_BIND_PASSWORD"),
			BaseDN:       required("LDAP_BASE_DN"),
		},
		Redis: RedisConfig{
			URL:      required("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		SMTP: SMTPConfig{
			Host:     required("SMTP_HOST"),
			Port:     smtpPort,
			User:     required("SMTP_USER"),
			Password: required("SMTP_PASSWORD"),
			From:     required("SMTP_FROM"),
		},
		Ramtun: RamtunConfig{
			Host:   required("RAMTUN_HOST"),
			ApiKey: required("RAMTUN_API_KEY"),
		},
	}

	if len(missing) > 0 {
		return nil, errors.New("variables de entorno requeridas no definidas: " + strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
