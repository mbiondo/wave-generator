# 🌊 Wave Generator

**Wave Generator** is a Go application that transforms images into smooth mathematical representations using piecewise cubic polynomials and renders them as scalable SVGs. Ideal for exploring the mathematical essence of skylines, waveforms, or any visually wave-like pattern.

[![Go CI](https://github.com/mbiondo/wave-generator/actions/workflows/ci.yml/badge.svg)](https://github.com/mbiondo/wave-generator/actions)

---

## 🎯 What does this project do?

Starting from an image input (like a skyline), it performs:

1. 🖤 Grayscale conversion
2. 🧠 Edge detection via vertical gradient analysis
3. 🧮 Curve fitting with least squares and cubic polynomials
4. 🖼️ SVG rendering and math output

You can use it via [API](#api) or [Playground (coming soon)](#).

---

## 🧪 Example

Imagine uploading this image:

![placeholder](./static/example-original.jpg)

And getting this as the resulting wave pattern:

![placeholder](./static/example-wave.svg)

---

## 🚀 Quick Install

```bash
# Clone the repo
git clone https://github.com/mbiondo/wave-generator.git
cd wave-generator

# Build
make build

# Or run with Docker
docker compose up --build
```

---

## 🔌 Friendly API

### Get an API Key

```bash
curl -X POST http://localhost:1155/generate-apikey
```

### Generate wave from image

```bash
curl -X POST http://localhost:1155/generate-wave \
     -H "X-API-Key: your_api_key" \
     -H "Content-Type: image/jpeg" \
     --data-binary "@./your-image.jpg"
```

📎 [View full API documentation →](./docs/api-docs.md)

---

## 📚 Key Resources

* 📘 [Math & Code Tutorial](./blog/wave-generator-math-tutorial.md)
* 📡 [API Documentation](./docs/api-docs.md)
* 🧾 [OpenAPI Spec](./docs/openapi.yaml)

---

## 📁 Project Structure

```
wave-generator/
├── handlers/       # HTTP endpoints
├── models/         # Data models
├── services/       # Core logic: image → wave → SVG
├── blog/           # Markdown tutorials
├── static/         # Static assets (images, SVGs)
├── docs/           # API & OpenAPI docs
├── Dockerfile      # Ready-to-run container
└── main.go         # Entry point
```

---

## 🧠 Motivation

This project began with a question: can we see an image as an equation? By applying math tools to visual inputs, we uncover beauty, structure, and creative possibilities.

---

## 🤝 Contributions

Pull requests welcome! Areas you can help with:

* New fitting techniques
* Improved edge detection strategies
* Optimizing for high-res images
* More output formats (PNG, PDF)

---

## 📄 License

MIT © Maximiliano Biondo
