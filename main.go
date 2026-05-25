package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Ошибка/Устаревание: r.URL.Path в старых версиях Go требовал проверки,
	// но в новом маршрутизаторе (Go 1.22+) точные совпадения работают строже.
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
	// Лучшая практика: Всегда создавайте явный и изолированный маршрутизатор (mux),
	// вместо использования глобального http.HandleFunc. Это защищает от скрытых багов.
	mux := http.NewServeMux()

	// В Go 1.22+ для точного совпадения корня (и только корня) используется синтаксис "GET /{$}"
	// Если оставить просто "/", он будет срабатывать на любые пути (например, /abc/def).
	mux.HandleFunc("GET /{$}", handler)

	// Лучшая практика: Всегда создавайте http.Server вручную и задавайте таймауты.
	// Дефолтный http.ListenAndServe без таймаутов уязвим к зависанию соединений (Slowloris атаки).
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Лучшая практика: Используйте пакет "log" вместо fmt.Println для логов сервера.
	log.Println("Starting server on :8080...")

	// Ошибка: Использование panic(err) завершает программу некорректно.
	// http.ErrServerClosed — это нормальное поведение при остановке, его не нужно считать ошибкой.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %s", err)
	}
}
