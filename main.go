package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Define handler functions
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/health", healthHandler)

    // Start server
    port := ":8080"
    fmt.Printf("Server is running on http://localhost%s\n", port)
    fmt.Println("Press Ctrl+C to stop the server")
    
    err := http.ListenAndServe(port, nil)
    if err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}

// homeHandler handles requests to root path
func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to my first Go web server!")
}

// helloHandler returns a personalized greeting
func helloHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "Guest"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

// healthHandler returns server health status
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "ok", "message": "Server is running properly"}`)
}
