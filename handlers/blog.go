package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

func BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	mdPath := filepath.Join("blog", "wave-generator-math-tutorial.md")
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}
	var htmlBuf []byte
	htmlBuf, err = renderMarkdown(mdBytes)
	if err != nil {
		http.Error(w, "Error rendering markdown", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<html><head><title>Wave Generator Blog</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
		<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
		<script src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 2rem auto; background: #f8fafd; color: #222; }
			pre, code { background: #f4f4f4; border-radius: 4px; }
			h1,h2,h3 { color: #3498db; }
			a { color: #2980b9; }
			.hljs { background: #f4f4f4; }
		</style>
	</head><body>
	<script>hljs.highlightAll();</script>
	`))
	w.Write(htmlBuf)
	w.Write([]byte(`</body></html>`))
}

func renderMarkdown(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
