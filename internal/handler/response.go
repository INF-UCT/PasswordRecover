package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type response struct {
	Success bool   `json:"success"`
	Meta    meta   `json:"meta"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
}

type errorBody struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

const errorBaseURL = "https://chpass.inf.uct.cl/api/errors/"

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeSuccess(w http.ResponseWriter, r *http.Request, data any) {
	writeJSON(w, http.StatusOK, response{
		Success: true,
		Meta: meta{
			RequestID: generateUUID(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Data:   data,
		Errors: nil,
	})
}

func writeError(w http.ResponseWriter, r *http.Request, status int, errType, title, detail string) {
	writeJSON(w, status, errorBody{
		Type:     errorBaseURL + errType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	})
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
