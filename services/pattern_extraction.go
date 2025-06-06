package services

import (
	"image"
	"math"
)

// ExtractPattern extracts a wave pattern from a grayscale image by detecting
// the strongest vertical gradient for each column in the image.
//
// The function scans each column (x-coordinate) of the image and finds the pixel position
// where there is the maximum difference in brightness between adjacent vertical pixels.
// This position is considered to be part of the wave pattern.
//
// Parameters:
//   - gray: A pointer to an image.Gray representing the grayscale input image
//   - w: The width of the image in pixels
//   - h: The height of the image in pixels
//
// Returns:
//   - A slice of float64 values with length w, where each value represents the y-coordinate
//     of the pattern at the corresponding x-coordinate
func ExtractPattern(gray *image.Gray, w, h int) []float64 {
	pattern := make([]float64, w)

	for x := range w {
		maxGradient := 0.0
		maxY := 0
		sumGradient := 0.0
		count := 0

		for y := 1; y < h; y++ {
			gradient := math.Abs(float64(gray.GrayAt(x, y).Y) - float64(gray.GrayAt(x, y-1).Y))
			sumGradient += gradient
			count++

			if gradient > maxGradient {
				maxGradient = gradient
				maxY = y
			}
		}

		// Use a weighted average to smooth the pattern extraction
		averageGradient := sumGradient / float64(count)
		if maxGradient > averageGradient*1.5 { // Threshold to avoid noise
			pattern[x] = float64(maxY)
		} else {
			pattern[x] = float64(h) / 2 // Default to middle of the image if no strong gradient
		}
	}

	return pattern
}
