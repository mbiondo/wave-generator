package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wave-generator/handlers"
)

func setupHandlers(mux *http.ServeMux) error {
	if mux == nil {
		return fmt.Errorf("nil ServeMux provided")
	}

	// Serve static files from the static directory
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve docs as static markdown rendered HTML
	mux.HandleFunc("/docs/api-docs", logHandler(handlers.APIDocsHandler))

	// Blog post handler
	mux.HandleFunc("/blog/wave-generator-math-tutorial", logHandler(handlers.BlogPostHandler))

	// API endpoints
	mux.HandleFunc("/generate-wave", logHandler(handlers.WavePatternHandler))
	mux.HandleFunc("/generate-apikey", logHandler(handlers.GenerateAPIKeyHandler))

	// Root handler must be last
	mux.HandleFunc("/", logHandler(handlers.IndexHandler))

	return nil
}

func logHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func startServer(mux *http.ServeMux) error {
	if mux == nil {
		return fmt.Errorf("nil ServeMux provided")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8899"
	}
	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

func main() {
	mux := http.NewServeMux()
	if err := setupHandlers(mux); err != nil {
		log.Fatal("Failed to setup handlers:", err)
	}
	if err := startServer(mux); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
