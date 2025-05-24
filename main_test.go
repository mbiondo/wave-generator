package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {

	// Test server setup
	t.Run("server setup", func(t *testing.T) {
		mux := http.NewServeMux()
		if err := setupHandlers(mux); err != nil {
			t.Errorf("setupHandlers failed: %v", err)
		}

		// Test static file serving
		req := httptest.NewRequest("GET", "/static/index.html", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		// A 301 redirect is expected for directories without trailing slash, but for files, accept 200 or 301
		if w.Code != http.StatusOK && w.Code != http.StatusMovedPermanently {
			t.Errorf("expected status OK or 301 for static file, got %v", w.Code)
		}

		// Test root handler
		req = httptest.NewRequest("GET", "/", nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK for root, got %v", w.Code)
		}

		// Test API endpoint
		req = httptest.NewRequest("POST", "/generate-wave", nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Error("expected error for empty POST to /generate-wave")
		}
	})

	// Test server startup with different PORT configurations
	t.Run("server port configuration", func(t *testing.T) {
		tests := []struct {
			name    string
			port    string
			wantErr bool
		}{
			{"default port", "", false},
			{"custom port", "8080", false},
			{"invalid port", "invalid", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.port != "" {
					if err := os.Setenv("PORT", tt.port); err != nil {
						t.Fatalf("Setenv failed: %v", err)
					}
				} else {
					if err := os.Unsetenv("PORT"); err != nil {
						t.Fatalf("Unsetenv failed: %v", err)
					}
				}
				origPort := os.Getenv("PORT")
				defer func() { _ = os.Setenv("PORT", origPort) }()

				mux := http.NewServeMux()
				err := setupHandlers(mux)
				if err != nil {
					t.Fatal(err)
				}

				serverErr := make(chan error, 1)
				server := &http.Server{
					Addr:    ":" + tt.port,
					Handler: mux,
				}

				go func() {
					serverErr <- server.ListenAndServe()
				}()

				// Give server time to start
				time.Sleep(100 * time.Millisecond)

				// Shutdown server
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					t.Logf("server shutdown: %v", err)
				}

				err = <-serverErr
				// Accept http.ErrServerClosed as a normal shutdown (not a test failure)
				if err != nil && err.Error() == "http: Server closed" {
					err = nil
				}
				if (err != nil && !tt.wantErr) || (err == nil && tt.wantErr) {
					t.Errorf("startServer() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}

func TestSetupHandlers(t *testing.T) {
	tests := []struct {
		name    string
		mux     *http.ServeMux
		wantErr bool
	}{
		{"valid mux", http.NewServeMux(), false},
		{"nil mux", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setupHandlers(tt.mux)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupHandlers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
