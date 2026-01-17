package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type BotHealth struct {
	Status    string `json:"status"`
	Bot       string `json:"bot"`
	Connected bool   `json:"connected"`
	Uptime    string `json:"uptime,omitempty"`
	Error     string `json:"error,omitempty"`
}

type AggregatedHealth struct {
	Status     string      `json:"status"`
	AllUp      bool        `json:"all_up"`
	TotalBots  int         `json:"total_bots"`
	BotsUp     int         `json:"bots_up"`
	BotsDown   int         `json:"bots_down"`
	Bots       []BotHealth `json:"bots"`
	CheckedAt  string      `json:"checked_at"`
}

var botEndpoints = []struct {
	Name string
	URL  string
}{
	{"eclub-bot", "http://localhost:8080/health"},
	{"visibleseed-bot", "http://localhost:8081/health"},
}

func checkBot(name, url string) BotHealth {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return BotHealth{
			Status:    "error",
			Bot:       name,
			Connected: false,
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()

	var health BotHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return BotHealth{
			Status:    "error",
			Bot:       name,
			Connected: false,
			Error:     "failed to decode response",
		}
	}

	return health
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var wg sync.WaitGroup
	results := make([]BotHealth, len(botEndpoints))

	for i, bot := range botEndpoints {
		wg.Add(1)
		go func(idx int, name, url string) {
			defer wg.Done()
			results[idx] = checkBot(name, url)
		}(i, bot.Name, bot.URL)
	}

	wg.Wait()

	botsUp := 0
	for _, r := range results {
		if r.Status == "ok" && r.Connected {
			botsUp++
		}
	}

	allUp := botsUp == len(botEndpoints)
	status := "ok"
	if !allUp {
		status = "degraded"
		if botsUp == 0 {
			status = "error"
		}
	}

	response := AggregatedHealth{
		Status:    status,
		AllUp:     allUp,
		TotalBots: len(botEndpoints),
		BotsUp:    botsUp,
		BotsDown:  len(botEndpoints) - botsUp,
		Bots:      results,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if !allUp {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"service": "discord-bots-health-aggregator"})
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Health aggregator listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
