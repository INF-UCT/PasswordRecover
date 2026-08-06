package ramtun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chpassword/internal/config"
)

const (
	ldapPasswordSyncPath = "/api/users/ldap-password-sync"
	maxAttempts          = 3
	requestTimeout       = 10 * time.Second
	retryBackoffBase     = 200 * time.Millisecond
)

type Client struct {
	cfg  config.RamtunConfig
	http *http.Client
}

type updatePasswordDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type statusError struct {
	code int
}

func (e statusError) Error() string {
	return fmt.Sprintf("ramtun: status %d", e.code)
}

func New(cfg config.RamtunConfig) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: requestTimeout},
	}
}

// SyncPassword notifica a ramtun del cambio de contraseña vía PATCH.
// Reintenta hasta maxAttempts ante errores de red o respuestas 5xx;
// no reintenta ante respuestas 4xx.
func (c *Client) SyncPassword(ctx context.Context, email, password string) error {
	body, err := json.Marshal(updatePasswordDto{Email: email, Password: password})
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := c.do(ctx, body)
		if err == nil {
			return nil
		}
		lastErr = err

		if !isRetryable(err) {
			return err
		}
		if attempt < maxAttempts {
			select {
			case <-time.After(retryBackoffBase * time.Duration(1<<(attempt-1))):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return lastErr
}

func (c *Client) do(ctx context.Context, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.cfg.Host+ldapPasswordSyncPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.ApiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return statusError{code: resp.StatusCode}
}

func isRetryable(err error) bool {
	se, ok := err.(statusError)
	if ok {
		return se.code >= 500
	}
	return true
}
