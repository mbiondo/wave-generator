# Wave Generator API Documentation

Welcome to the **Wave Generator** public API documentation.

---

## Table of Contents

- [Wave Generator API Documentation](#wave-generator-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Authentication \& Rate Limiting](#authentication--rate-limiting)
  - [Endpoints](#endpoints)
    - [Generate API Key](#generate-api-key)
    - [Generate Wave](#generate-wave)
  - [Response Format](#response-format)
  - [Error Handling](#error-handling)
  - [Examples](#examples)
  - [References](#references)

---

## Overview

The Wave Generator API allows you to upload an image and receive a mathematical and SVG representation of its wave pattern, segmented into cubic polynomials.

---

## Authentication & Rate Limiting

- **API Key Required:** All requests (except same-origin/browser) must include an API key in the `X-API-Key` header.
- **Rate Limit:** 1000 requests per hour per API key.
- **How to get an API key:** See [Generate API Key](#generate-api-key).

---

## Endpoints

### Generate API Key

**POST** `/generate-apikey`

Generate a new API key for authenticating requests.

**Request:**
```bash
curl -X POST http://localhost:1155/generate-apikey
```

**Response:**
```json
{
  "api_key": "api_1712345678901234567"
}
```

---

### Generate Wave

**POST** `/generate-wave`

Extracts the wave pattern from an image, fits cubic polynomial segments, and returns SVG and segment data.

**Headers:**
- `X-API-Key: <your-api-key>` (required unless same-origin)

**Body:**  
Image file (PNG or JPEG), sent as binary.

**Response:**
```json
{
  "segments": [
    {
      "domain_start": 0,
      "domain_end": 15,
      "a3": 0.0001,
      "a2": -0.02,
      "a1": 1.5,
      "a0": 10.0,
      "expression": "for x ∈ [0,15]: y = ...",
      "svg": "<svg ...>...</svg>"
    }
    // ...
  ],
  "svg": "<svg ...>...</svg>",
  "segment_svgs": ["<svg ...>", ...],
  "coords": [[0, 12], [1, 13], ...]
}
```

---

## Response Format

- `segments`: Array of polynomial segment objects, each with coefficients, domain, expression, and a mini SVG.
- `svg`: SVG string representing the full wave.
- `segment_svgs`: Array of SVG strings, one per segment.
- `coords`: Array of `[x, y]` pairs representing the extracted pattern.

---

## Error Handling

- `401 Unauthorized`: Missing or invalid API key.
- `429 Too Many Requests`: Rate limit exceeded.
- `400 Bad Request`: Invalid image.
- `422 Unprocessable Entity`: Processing error.

---

## Examples

**Generate API Key**
```bash
curl -X POST http://localhost:1155/generate-apikey
```

**Generate Wave**
```bash
curl -X POST http://localhost:1155/generate-wave \
     -H "X-API-Key: api_1712345678901234567" \
     -H "Content-Type: image/jpeg" \
     --data-binary "@./image.jpg"
```

---

## References

- [OpenAPI YAML (raw)](https://github.com/mbiondo/wave-generator/blob/main/docs/openapi.yaml)
- [Math & Programming Tutorial](../blog/wave-generator-math-tutorial.md)
