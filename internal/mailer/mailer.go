package mailer

import (
	"fmt"
	"log/slog"

	"chpassword/internal/config"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	cfg config.SMTPConfig
}

func New(cfg config.SMTPConfig) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendResetEmail(to, resetURL string) error {
	slog.Info("preparando correo de reset",
		"smtp_host", m.cfg.Host,
		"smtp_port", m.cfg.Port,
		"smtp_user", m.cfg.User,
		"from", m.cfg.From,
		"to", to,
	)

	msg := mail.NewMsg()

	if err := msg.From(m.cfg.From); err != nil {
		slog.Error("dirección remitente inválida", "from", m.cfg.From, "error", err)
		return fmt.Errorf("dirección remitente inválida: %w", err)
	}
	if err := msg.To(to); err != nil {
		slog.Error("dirección destinatario inválida", "to", to, "error", err)
		return fmt.Errorf("dirección destinatario inválida: %w", err)
	}

	msg.Subject("Restablecimiento de contraseña")
	msg.SetBodyString(mail.TypeTextHTML, fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<body>
  <p>Hola,</p>
  <p>Recibimos una solicitud para restablecer tu contraseña.</p>
  <p>Haz clic en el siguiente enlace para continuar:</p>
  <p><a href="%s">%s</a></p>
  <p>Este enlace expira en 10 minutos.</p>
  <p>Si no solicitaste este cambio, puedes ignorar este correo.</p>
</body>
</html>`, resetURL, resetURL))

	slog.Info("conectando al servidor SMTP", "host", m.cfg.Host, "port", m.cfg.Port)

	client, err := mail.NewClient(
		m.cfg.Host,
		mail.WithPort(m.cfg.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(m.cfg.User),
		mail.WithPassword(m.cfg.Password),
		mail.WithSSL(),
	)
	if err != nil {
		slog.Error("error creando cliente SMTP", "host", m.cfg.Host, "port", m.cfg.Port, "error", err)
		return err
	}

	slog.Info("enviando correo", "to", to)
	if err := client.DialAndSend(msg); err != nil {
		slog.Error("error en DialAndSend", "to", to, "host", m.cfg.Host, "port", m.cfg.Port, "error", err)
		return err
	}

	slog.Info("correo enviado exitosamente", "to", to)
	return nil
}
