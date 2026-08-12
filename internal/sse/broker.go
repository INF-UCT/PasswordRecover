package sse

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type PasswordChangeEvent struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

type Broker struct {
	Notifier       chan PasswordChangeEvent
	newClients     chan chan PasswordChangeEvent
	closingClients chan chan PasswordChangeEvent
	clients        map[chan PasswordChangeEvent]bool
}

func NewBroker() *Broker {
	return &Broker{
		Notifier:       make(chan PasswordChangeEvent, 10),
		newClients:     make(chan chan PasswordChangeEvent),
		closingClients: make(chan chan PasswordChangeEvent),
		clients:        make(map[chan PasswordChangeEvent]bool),
	}
}

func (broker *Broker) Listen() {
	for {
		select {
		case subscriber := <-broker.newClients:
			broker.clients[subscriber] = true
			slog.Info("[SSE] New client connected")

		case subscriber := <-broker.closingClients:
			delete(broker.clients, subscriber)
			slog.Info("[SSE] Client disconnected")

		case event := <-broker.Notifier:
			for clientMessageChan := range broker.clients {
				select {
				case clientMessageChan <- event:
					slog.Info("[SSE] Message send")
				default:
					slog.Warn("[SSE] Skipping event")
				}
			}
		}
	}
}

func (broker *Broker) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "-", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("X-Accel-Buffering", "no")

	messageChan := make(chan PasswordChangeEvent, 10)
	broker.newClients <- messageChan
	defer func() { broker.closingClients <- messageChan }()

	for {
		select {
		case <-request.Context().Done():
			return
		case event := <-messageChan:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
