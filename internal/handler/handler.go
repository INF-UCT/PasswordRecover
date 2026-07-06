package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"chpassword/internal/ldap"
	"chpassword/internal/mailer"
	"chpassword/internal/redis"
)

type Handler struct {
	ldap        *ldap.Client
	redis       *redis.Client
	mailer      *mailer.Mailer
	baseURL     string
	frontendURL string
}

func New(l *ldap.Client, r *redis.Client, m *mailer.Mailer, baseURL, frontendURL string) *Handler {
	return &Handler{
		ldap:        l,
		redis:       r,
		mailer:      m,
		baseURL:     baseURL,
		frontendURL: frontendURL,
	}
}

// RequestReset recibe un correo y envía un enlace de restablecimiento si el usuario existe en LDAP.
func (h *Handler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" {
		writeError(w, r, http.StatusBadRequest, "invalid-input", "Datos inválidos", "El campo email es requerido")
		return
	}

	slog.Info("password reset request", "email", body.Email)

	exists, _, err := h.ldap.UserExists(body.Email)
	if err != nil {
		slog.Error("error consultando LDAP", "email", body.Email, "error", err)
		writeError(w, r, http.StatusInternalServerError, "internal", "Error interno", "No se pudo procesar la solicitud")
		return
	}

	if exists {
		token := generateUUID()
		slog.Info("usuario encontrado en LDAP, generando token", "email", body.Email, "token", token)

		if err := h.redis.SaveToken(r.Context(), token, body.Email); err != nil {
			slog.Error("error guardando token en Redis", "email", body.Email, "error", err)
			writeError(w, r, http.StatusInternalServerError, "internal", "Error interno", "No se pudo procesar la solicitud")
			return
		}

		resetURL := h.frontendURL + "/reset/" + token

		if err := h.mailer.SendResetEmail(body.Email, resetURL); err != nil {
			slog.Error("error enviando correo", "email", body.Email, "error", err)
		} else {
			slog.Info("correo de reset enviado", "email", body.Email)
		}
	} else {
		slog.Info("email no encontrado en LDAP, respondiendo igual para evitar enumeración", "email", body.Email)
	}

	// Siempre responde 200 para no revelar si el correo existe o no.
	writeSuccess(w, r, map[string]string{
		"message": "Si el correo existe en el sistema, recibirás un enlace para cambiar tu contraseña.",
	})
}

// ValidateToken verifica si un token sigue siendo válido.
func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	slog.Info("validando token", "token", token)

	exists, err := h.redis.TokenExists(r.Context(), token)
	if err != nil {
		slog.Error("error consultando Redis", "token", token, "error", err)
		writeError(w, r, http.StatusInternalServerError, "internal", "Error interno", "No se pudo verificar el token")
		return
	}

	if !exists {
		slog.Warn("token no encontrado", "token", token)
		writeError(w, r, http.StatusNotFound, "not-found", "Token no encontrado", "El token no existe o ya fue utilizado")
		return
	}

	ttl, err := h.redis.GetTokenTTL(r.Context(), token)
	if err != nil || ttl <= 0 {
		slog.Warn("token vencido", "token", token)
		writeError(w, r, http.StatusGone, "token-expired", "Token vencido", "El token ha expirado")
		return
	}
	slog.Info("token válido", "token", token, "ttl", ttl)

	writeSuccess(w, r, map[string]any{
		"valid":      true,
		"expires_at": time.Now().UTC().Add(ttl).Format(time.RFC3339),
	})
}

// ConfirmReset valida el token y aplica el cambio de contraseña en LDAP.
func (h *Handler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token             string `json:"token"`
		NewPassword       string `json:"new_password"`
		ConfirmedPassword string `json:"confirmed_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid-input", "Datos inválidos", "El cuerpo de la solicitud es inválido")
		return
	}

	if body.Token == "" {
		writeValidationError(w, r, []fieldError{{
			Field:   "token",
			Code:    "required",
			Message: "El campo token es requerido.",
		}})
		return
	}

	if errs := validatePassword(body.NewPassword); len(errs) > 0 {
		writeValidationError(w, r, errs)
		return
	}

	if mismatch := validatePasswordMatch(body.NewPassword, body.ConfirmedPassword); mismatch != nil {
		writeValidationError(w, r, []fieldError{*mismatch})
		return
	}

	slog.Info("confirmando cambio de contraseña", "token", body.Token)

	email, err := h.redis.GetEmailByToken(r.Context(), body.Token)
	if err != nil {
		slog.Warn("token no encontrado al confirmar", "token", body.Token, "error", err)
		writeError(w, r, http.StatusNotFound, "not-found", "Token no encontrado", "El token no existe o ya fue utilizado")
		return
	}

	_, userDN, err := h.ldap.UserExists(email)
	if err != nil || userDN == "" {
		slog.Error("error consultando LDAP al confirmar", "email", email, "error", err)
		writeError(w, r, http.StatusInternalServerError, "internal", "Error interno", "No se pudo verificar el usuario")
		return
	}

	if err := h.ldap.ChangePassword(userDN, body.NewPassword); err != nil {
		if errors.Is(err, ldap.ErrPasswordPolicy) {
			slog.Warn("contraseña rechazada por política LDAP", "email", email, "error", err)
			writeError(w, r, http.StatusUnprocessableEntity, "password-policy", "Contraseña inválida", err.Error())
			return
		}
		slog.Error("error cambiando contraseña en LDAP", "email", email, "error", err)
		writeError(w, r, http.StatusInternalServerError, "internal", "Error interno", "No se pudo actualizar la contraseña")
		return
	}

	slog.Info("contraseña actualizada correctamente", "email", email)

	if err := h.redis.DeleteToken(r.Context(), body.Token); err != nil {
		slog.Warn("no se pudo eliminar el token tras cambio exitoso", "token", body.Token, "error", err)
	}

	writeSuccess(w, r, map[string]string{
		"message": "Contraseña actualizada correctamente.",
	})
}
