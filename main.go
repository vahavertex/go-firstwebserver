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
	"strconv"
	"syscall"
	"time"
)

type PageData struct {
	Path string
}

type Config struct {
	Port string
}

// Выносим шаблон на уровень пакета (глобальная переменная для кэша)
var tmpl *template.Template

func LoadConfig() (Config, error) {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	portPtr := flag.String("port", envPort, "Port to listen on (1-65355)")
	flag.Parse()

	// Логическая проверка: валиден ли порт вообще?
	p, err := strconv.Atoi(*portPtr)
	if err != nil || p < 1 || p > 65535 {
		return Config{}, fmt.Errorf("invalid port number: %s", *portPtr)
	}

	return Config{Port: *portPtr}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Path: r.URL.Path}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// ВЫСОКАЯ ПРОИЗВОДИТЕЛЬНОСТЬ: Теперь мы берем уже готовый шаблон из памяти.
	// Никакого обращения к диску! Операция стала в тысячи раз быстрее.
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Инициализируем (кэшируем) шаблон строго ОДИН раз при старте.
	// Если файла нет — программа упадет сразу при запуске, а не во время работы.
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Critical: failed to parse template: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Канал для отслеживания критических ошибок самого сервера во время старта
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting server. Open in browser: http://localhost:%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// ИСПРАВЛЕНИЕ ЛОГИКИ: Теперь main ждет ЛИБО сигнал от ОС, ЛИБО ошибку старта сервера
	select {
	case err := <-serverErrors:
		log.Fatalf("Critical: Server failed to start: %v", err)
	case <-ctx.Done():
		log.Println("\nShutdown signal received. Closing server gracefully...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server stopped. Port %s is free.\n", cfg.Port)
}
