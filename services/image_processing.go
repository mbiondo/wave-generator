package services

import (
	"image"
)

// ToGray converts any image.Image to a grayscale image.
// It creates a new grayscale image with the same dimensions as the input image
// and copies each pixel, automatically converting the colors to grayscale values.
// Returns a pointer to the new grayscale image.
func ToGray(img image.Image) *image.Gray {
	b := img.Bounds()
	gray := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}
