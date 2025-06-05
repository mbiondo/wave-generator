package services

import (
	"image"
	"image/color"
	"math"
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

// ToGraySmooth converts any image.Image to a grayscale image with noise reduction.
// It applies a simple box blur to smooth the image and reduce noise before converting to grayscale.
// Returns a pointer to the new grayscale image.
func ToGraySmooth(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	// Apply a simple box blur for noise reduction
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, count float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						c := img.At(px, py)
						rr, gg, bb, _ := c.RGBA()
						r += float64(rr)
						g += float64(gg)
						b += float64(bb)
						count++
					}
				}
			}
			avgR := uint8(r / count / 256)
			avgG := uint8(g / count / 256)
			avgB := uint8(b / count / 256)
			gray.Set(x, y, color.Gray{Y: uint8((avgR + avgG + avgB) / 3)})
		}
	}

	return gray
}

// EdgeDetection applies the Sobel filter to detect edges in an image.
// It highlights the silhouette of objects like skylines or landscapes.
// Returns a pointer to the new grayscale image with edges highlighted.
func EdgeDetection(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	// Sobel kernels
	sobelX := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var gx, gy int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						c := img.At(px, py)
						r, g, b, _ := c.RGBA()
						grayValue := int((r + g + b) / 3 / 256)
						gx += grayValue * sobelX[ky+1][kx+1]
						gy += grayValue * sobelY[ky+1][kx+1]
					}
				}
			}
			// Calculate gradient magnitude
			magnitude := uint8(clamp(int(math.Sqrt(float64(gx*gx+gy*gy))), 0, 255))
			gray.Set(x, y, color.Gray{Y: magnitude})
		}
	}

	return gray
}

// clamp limits a value to a specified range.
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
