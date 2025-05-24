package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to create a minimal markdown file for testing
func createTempMarkdown(t *testing.T, path, content string) {
	t.Helper()
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func TestBlogPostHandler_OK(t *testing.T) {
	tmpPath := "blog/wave-generator-math-tutorial.md"
	createTempMarkdown(t, tmpPath, "# Hello\nThis is a test blog.")
	defer func() {
		_ = os.Remove(tmpPath)
		_ = os.Remove("static/page_template.html")
		_ = os.Remove("static")
		_ = os.Remove("blog")
	}()

	// Ensure template exists for the test
	templatePath := "static/page_template.html"
	templateContent := `
	<!DOCTYPE html>
	<html><head><title>{{TITLE}}</title></head>
	<body>
	<header class="header-nav"></header>
	<!--CONTENT-->
	</body></html>`
	_ = os.MkdirAll("static", 0755)
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("write template failed: %v", err)
	}
	defer func() { _ = os.Remove(templatePath) }()

	req := httptest.NewRequest("GET", "/blog/wave-generator-math-tutorial", nil)
	w := httptest.NewRecorder()
	BlogPostHandler(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "This is a test blog") {
		t.Errorf("response missing blog content")
	}
	if !strings.Contains(string(body), "<h1") {
		t.Errorf("response missing h1")
	}
}

func TestBlogPostHandler_NotFound(t *testing.T) {
	tmpPath := "blog/wave-generator-math-tutorial.md"
	_ = os.Remove(tmpPath) // ensure file does not exist

	req := httptest.NewRequest("GET", "/blog/wave-generator-math-tutorial", nil)
	w := httptest.NewRecorder()
	BlogPostHandler(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestAPIDocsHandler_OK(t *testing.T) {
	tmpPath := "docs/api-docs.md"
	createTempMarkdown(t, tmpPath, "# API Docs\nThis is the API documentation.")
	defer func() {
		_ = os.Remove(tmpPath)
		_ = os.Remove("static/page_template.html")
		_ = os.Remove("static")
		_ = os.Remove("docs")
	}()

	// Ensure template exists for the test
	templatePath := "static/page_template.html"
	templateContent := `
	<!DOCTYPE html>
	<html><head><title>{{TITLE}}</title></head>
	<body>
	<header class="header-nav"></header>
	<!--CONTENT-->
	</body></html>`
	_ = os.MkdirAll("static", 0755)
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("write template failed: %v", err)
	}
	defer func() { _ = os.Remove(templatePath) }()

	req := httptest.NewRequest("GET", "/docs/api-docs", nil)
	w := httptest.NewRecorder()
	APIDocsHandler(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "API documentation") {
		t.Errorf("response missing docs content")
	}
}

func TestAPIDocsHandler_NotFound(t *testing.T) {
	tmpPath := "docs/api-docs.md"
	_ = os.Remove(tmpPath) // ensure file does not exist

	req := httptest.NewRequest("GET", "/docs/api-docs", nil)
	w := httptest.NewRecorder()
	APIDocsHandler(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
