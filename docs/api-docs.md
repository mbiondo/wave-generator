# 📡 Wave Generator API Documentation

Welcome to the **Wave Generator** public API! This documentation will help you send images, receive fitted wave equations, and render mathematical beauty from pixel data.

---

## 🔑 Authentication & Rate Limiting

* All requests (except from the browser UI) must include an API key.
* Include your key in the header:

  ```
  X-API-Key: your_api_key_here
  ```
* **Rate limit**: 1000 requests per hour per key

---

## 🧪 How to Get an API Key

**POST** `/generate-apikey`

```bash
curl -X POST http://localhost:1155/generate-apikey
```

Response:

```json
{
  "api_key": "api_1712345678901234567"
}
```

---

## 🎯 Main Endpoint: `/generate-wave`

Submit an image (PNG or JPEG) and receive:

* SVG with full wave
* Polynomial equations (per segment)
* Segment mini SVGs
* Raw coordinate points

### Request Example

```bash
curl -X POST http://localhost:1155/generate-wave \
     -H "X-API-Key: api_..." \
     -H "Content-Type: image/png" \
     --data-binary "@./your-image.png"
```

---

### Response Format

```json
{
  "segments": [
    {
      "domain_start": 0,
      "domain_end": 20,
      "a3": -0.0003,
      "a2": 0.02,
      "a1": -0.6,
      "a0": 15.0,
      "expression": "for x ∈ [0,20]: y = ...",
      "svg": "<svg>...</svg>"
    }
    // more segments...
  ],
  "svg": "<svg>...</svg>",
  "segment_svgs": ["<svg>...</svg>", ...],
  "coords": [[0, 12], [1, 13], ...]
}
```

---

## ⚠️ Error Handling

- **400**: Invalid or unsupported image
- **401**: API key missing/invalid
- **422**: Processing failed
- **429**: Rate limit exceeded

---

## 📘 Learn More

* [Math + Code Tutorial](../blog/wave-generator-math-tutorial)
* [API Reference YAML](https://github.com/mbiondo/wave-generator/blob/main/docs/openapi.yaml)
* [Playground](#coming-soon)

---

## 🎨 Use Cases

* Generate skyline curves for generative art
* Convert voice or waveform images into equations
* Build educational demos to visualize how polynomials fit data

If you build something awesome, open an issue or share your demo!
