package services

import (
	"image"
	"testing"
)

func TestExtractPattern(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 10, 10))

	pattern := ExtractPattern(img, 10, 10)

	if len(pattern) != 10 {
		t.Fatalf("expected pattern length 10, got %d", len(pattern))
	}

	for _, y := range pattern {
		if y < 0 || y >= 10 {
			t.Errorf("invalid pattern y-coordinate: %f", y)
		}
	}
}
