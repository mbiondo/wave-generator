package services

import (
	"image"
	"image/color"
	"testing"
)

func TestToGray(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 0, A: 255})
		}
	}

	gray := ToGray(img)

	if gray.Bounds() != img.Bounds() {
		t.Fatalf("expected bounds %v, got %v", img.Bounds(), gray.Bounds())
	}

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			grayColor := gray.GrayAt(x, y).Y
			// Using standard grayscale conversion formula: 0.299*R + 0.587*G + 0.114*B
			expected := uint8(0.299*float64(x) + 0.587*float64(y) + 0.114*0)
			if abs(int(grayColor)-int(expected)) > 1 {
				t.Errorf("at (%d, %d): expected gray value %d (±1), got %d", x, y, expected, grayColor)
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
