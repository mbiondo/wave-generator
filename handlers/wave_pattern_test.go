package handlers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"wave-generator/models"
)

func TestWavePatternHandler(t *testing.T) {
	// Create test image with more variation
	img := image.NewGray(image.Rect(0, 0, 33, 10))
	for x := 0; x < 33; x++ {
		for y := 0; y < 10; y++ {
			// Create a wave-like pattern to avoid singular matrices
			val := uint8(((x+y)%10)*25 + 5)
			img.SetGray(x, y, color.Gray{Y: val})
		}
	}

	tests := []struct {
		name         string
		method       string
		body         []byte
		contentType  string
		wantStatus   int
		wantResponse bool
	}{
		{
			name:         "valid image",
			method:       http.MethodPost,
			contentType:  "image/png",
			wantStatus:   http.StatusOK,
			wantResponse: true,
		},
		{
			name:         "wrong method",
			method:       http.MethodGet,
			wantStatus:   http.StatusMethodNotAllowed,
			wantResponse: false,
		},
		{
			name:         "invalid image",
			method:       http.MethodPost,
			body:         []byte("not an image"),
			contentType:  "image/png",
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			if tt.name == "valid image" {
				if err := png.Encode(&body, img); err != nil {
					t.Fatal(err)
				}
				tt.body = body.Bytes()
			}

			req := httptest.NewRequest(tt.method, "/generate-wave", bytes.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rec := httptest.NewRecorder()

			WavePatternHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantResponse {
				var response models.ResponsePayload
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(response.Segments) == 0 {
					t.Error("expected non-empty segments")
				}
				if response.SVG == "" {
					t.Error("expected non-empty SVG")
				}
			}
		})
	}

	// Test singular matrix case
	t.Run("singular matrix case", func(t *testing.T) {
		uniformImg := image.NewGray(image.Rect(0, 0, 8, 8))
		var buf bytes.Buffer
		png.Encode(&buf, uniformImg)

		req := httptest.NewRequest(http.MethodPost, "/generate-wave", &buf)
		req.Header.Set("Content-Type", "image/png")
		rec := httptest.NewRecorder()

		WavePatternHandler(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})
}
