package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"chpassword/internal/config"
	"chpassword/internal/handler"
	"chpassword/internal/ldap"
	"chpassword/internal/mailer"
	"chpassword/internal/middleware"
	"chpassword/internal/redis"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("error cargando configuración", "error", err)
		os.Exit(1)
	}

	redisClient, err := redis.New(cfg.Redis, cfg.App.TokenTTLMinutes)
	if err != nil {
		slog.Error("error conectando a Redis", "error", err)
		os.Exit(1)
	}

	ldapClient := ldap.New(cfg.LDAP)
	mailerClient := mailer.New(cfg.SMTP)

	h := handler.New(ldapClient, redisClient, mailerClient, cfg.App.BaseURL, cfg.App.FrontendURL)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/password-reset/request", h.RequestReset)
	mux.HandleFunc("GET /api/v1/password-reset/{token}", h.ValidateToken)
	mux.HandleFunc("POST /api/v1/password-reset/confirm", h.ConfirmReset)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	slog.Info("servidor iniciado", "addr", addr)

	if err := http.ListenAndServe(addr, middleware.Logging(mux)); err != nil {
		slog.Error("error iniciando servidor", "error", err)
		os.Exit(1)
	}
}
