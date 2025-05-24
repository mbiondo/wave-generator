package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

func renderPageTemplate(title string, content []byte) ([]byte, error) {
	templatePath := filepath.Join("static", "page_template.html")
	tpl, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	html := string(tpl)
	// Si el contenido no tiene <h1>, agrega uno pero con clases y estructura consistente
	if !strings.Contains(string(content), "<h1") {
		content = append([]byte(`<h1 class="page-title">`+title+`</h1>`), content...)
	}
	// Asegura que el contenido esté dentro de <main> para mantener el padding y ancho
	if !strings.Contains(html, "<main>") {
		html = strings.Replace(html, "<!--CONTENT-->", `<main>`+string(content)+`</main>`, 1)
	} else {
		html = strings.Replace(html, "<!--CONTENT-->", string(content), 1)
	}
	html = strings.Replace(html, "{{TITLE}}", title, 1)
	return []byte(html), nil
}

func BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join("blog", "wave-generator-math-tutorial.md")
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}
	var htmlBuf bytes.Buffer
	if err := goldmark.Convert(mdBytes, &htmlBuf); err != nil {
		http.Error(w, "Error rendering markdown", http.StatusInternalServerError)
		return
	}
	page, err := renderPageTemplate("Wave Generator Blog", htmlBuf.Bytes())
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}

func renderMarkdown(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func APIDocsHandler(w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join("docs", "api-docs.md")
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		http.Error(w, "API docs not found", http.StatusNotFound)
		return
	}
	var htmlBuf bytes.Buffer
	if err := goldmark.Convert(mdBytes, &htmlBuf); err != nil {
		http.Error(w, "Error rendering markdown", http.StatusInternalServerError)
		return
	}
	page, err := renderPageTemplate("Wave Generator API Docs", htmlBuf.Bytes())
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}
