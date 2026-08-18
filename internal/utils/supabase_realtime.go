package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"nalakarsa/internal/config"
)

type realtimeBroadcastPayload struct {
	Messages []realtimeBroadcastMessage `json:"messages"`
}

type realtimeBroadcastMessage struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

// BroadcastToSupabase sends a real-time event to a Supabase WebSocket channel asynchronously
func BroadcastToSupabase(topic string, event string, payload interface{}, cfg *config.Config) {
	if cfg == nil || cfg.SupabaseURL == "" || cfg.SupabaseServiceRoleKey == "" {
		return
	}

	go func() {
		endpoint := fmt.Sprintf("%s/realtime/v1/api/broadcast", strings.TrimRight(cfg.SupabaseURL, "/"))

		bodyData := realtimeBroadcastPayload{
			Messages: []realtimeBroadcastMessage{
				{
					Topic:   topic,
					Event:   event,
					Payload: payload,
				},
			},
		}

		jsonData, err := json.Marshal(bodyData)
		if err != nil {
			log.Printf("[Supabase Realtime] Error marshalling payload: %v\n", err)
			return
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("[Supabase Realtime] Error creating request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", cfg.SupabaseServiceRoleKey)
		req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Supabase Realtime] Error sending broadcast: %v\n", err)
			return
		}
		defer resp.Body.Close()
	}()
}
