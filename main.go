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

// Config хранит настройки нашего сервера
type Config struct {
	Port string
}

// LoadConfig собирает настройки из флагов и переменных окружения
func LoadConfig() Config {
	// 1. Проверяем, есть ли переменная окружения PORT (например, в Docker)
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080" // Дефолтное значение
	}

	// 2. Добавляем флаг командной строки "-port"
	// Если флаг не передан в консоли, подставится значение из envPort
	portPtr := flag.String("port", envPort, "Port to listen on")

	// Обязательно вызываем Parse, чтобы Go прочитал аргументы из консоли
	flag.Parse()

	return Config{
		Port: *portPtr,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	data := PageData{Path: r.URL.Path}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func main() {
	// Загружаем конфигурацию при старте приложения
	cfg := LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler)

	// Форматируем адрес строки (добавляем двоеточие перед портом, если его нет)
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		// Динамически выводим порт в лог
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
