package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create temporary test directory structure
	tmpDir := filepath.Join(os.TempDir(), "wave-generator-test")
	staticDir := filepath.Join(tmpDir, "static")
	err := os.MkdirAll(staticDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test index.html in static directory
	indexPath := filepath.Join(staticDir, "index.html")
	content := "<html><body>Test</body></html>"
	if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test index.html: %v", err)
	}

	// Set working directory to test directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	// Test root path
	t.Run("root path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		IndexHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		if contentType := rec.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
			t.Errorf("expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
		}
	})

	// Test missing file
	t.Run("missing file", func(t *testing.T) {
		// Remove static directory to simulate missing file
		if err := os.RemoveAll(staticDir); err != nil {
			t.Fatalf("failed to remove static directory: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		IndexHandler(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	// Test invalid path
	t.Run("invalid path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		rec := httptest.NewRecorder()
		IndexHandler(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})
}
