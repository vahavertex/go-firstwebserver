package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/info", infoHandler)
	mux.HandleFunc("/api/time", timeHandler)

	// Custom server configuration
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	log.Println("Available endpoints:")
	log.Println("  GET /api/info - Server information")
	log.Println("  GET /api/time  - Current server time")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "This is my first Go web server!",
		Time:    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"current_time": time.Now().Format("15:04:05"),
		"date":         time.Now().Format("2006-01-02"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
