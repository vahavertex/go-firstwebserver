package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PageData struct {
	Path string
}

type Config struct {
	Port string
}

// Кэшируем шаблон при старте приложения.
// template.Must вызовет панику, если в HTML файле будет синтаксическая ошибка.
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func LoadConfig() Config {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	portPtr := flag.String("port", envPort, "Port to listen on")
	flag.Parse()

	return Config{
		Port: *portPtr,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Path: r.URL.Path}

	// Заголовок Content-Type лучше ставить ДО отправки любого статуса или данных
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Просто исполняем уже готовый в памяти шаблон
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		// Если шаблон упал в процессе передачи, заголовок ответа изменить уже нельзя,
		// но мы хотя бы логируем это событие.
	}
}

func main() {
	cfg := LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	// GET /{$} строго соответствует только корню "/"
	// Если нужно обрабатывать вообще все пути, измените на "GET /"
	mux.HandleFunc("GET /{$}", handler)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting server. Open in browser: http://localhost:%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	<-ctx.Done()
	log.Println("\nShutdown signal received. Closing server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server stopped. Port %s is free.\n", cfg.Port)
}
