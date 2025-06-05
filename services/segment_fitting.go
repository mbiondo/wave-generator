package services

import (
	"fmt"

	"wave-generator/models"

	"gonum.org/v1/gonum/mat"
)

// FitSegments takes a pattern array and width to fit cubic polynomials to segments of the pattern.
//
// This function divides the input pattern into a fixed number of segments (currently 3)
// and fits a cubic polynomial (form: ax³ + bx² + cx + d) to each segment using least squares fitting.
//
// Parameters:
//   - pattern: A float64 slice representing the y-values of the pattern
//   - width: The total width of the pattern (number of x-coordinates)
//
// Returns:
//   - A slice of PolySegment objects, each containing:
//   - The start and end x-coordinates of the segment
//   - The coefficients of the cubic polynomial (a₃, a₂, a₁, a₀)
//   - A string representation of the polynomial expression
//
// The function uses QR decomposition to solve the least squares problem for each segment.
// If the solver encounters an error or if a segment would have zero or negative width,
// the function will panic with an appropriate error message.
func FitSegments(pattern []float64, width int) []models.PolySegment {
	// Input validation
	if pattern == nil || width <= 0 || len(pattern) < 4 {
		panic("invalid input: pattern array must not be nil and width must be positive")
	}

	// Validate width against data length
	if width > len(pattern) {
		width = len(pattern)
	}

	// Minimum points needed for cubic polynomial fit
	const minPointsPerSeg = 4

	// Calculate maximum possible segments based on input size
	maxSeg := width / minPointsPerSeg

	// Calculate number of segments based on image width
	// For small images (width < 128), use width/16 segments
	// For larger images, use up to 32 segments
	nSeg := width / 16
	if nSeg > 32 {
		nSeg = 32
	}
	if nSeg < 1 {
		nSeg = 1
	}
	if maxSeg < nSeg {
		nSeg = maxSeg
	}

	// Ensure minimum segment width
	segW := width / nSeg
	if segW < minPointsPerSeg {
		nSeg = 1
		segW = width
	}

	segments := make([]models.PolySegment, 0, nSeg)

	for i := range nSeg {
		x0 := i * segW
		x1 := (i + 1) * segW
		if x0 >= width {
			break // Avoid out-of-bounds segments
		}
		if x1 > width {
			x1 = width
		}

		m := x1 - x0
		if m < minPointsPerSeg {
			continue // Skip segments with too few points
		}

		X := mat.NewDense(m, 4, nil)
		Y := mat.NewVecDense(m, nil)
		allSame := true
		firstY := pattern[x0]
		for j := 0; j < m; j++ {
			xVal := float64(x0 + j)
			X.Set(j, 0, xVal*xVal*xVal) // math.Pow(xVal, 3)
			X.Set(j, 1, xVal*xVal)      // math.Pow(xVal, 2)
			X.Set(j, 2, xVal)
			X.Set(j, 3, 1)
			Y.SetVec(j, pattern[x0+j])
			if pattern[x0+j] != firstY {
				allSame = false
			}
		}

		if allSame {
			continue // skip this segment, don't panic
		}

		var qr mat.QR
		qr.Factorize(X)

		var c mat.VecDense
		if err := qr.SolveVecTo(&c, false, Y); err != nil {
			panic(fmt.Sprintf("Error solving system: %v", err))
		}

		cv := c.RawVector().Data
		expr := fmt.Sprintf(
			"for x ∈ [%d,%d]: y = %.6f·x³ %+.6f·x² %+.6f·x %+.6f",
			x0, x1-1,
			cv[0], cv[1], cv[2], cv[3],
		)

		segments = append(segments, models.PolySegment{
			X0:         x0,
			X1:         x1 - 1,
			CoefA3:     cv[0],
			CoefA2:     cv[1],
			CoefA1:     cv[2],
			CoefA0:     cv[3],
			Expression: expr,
		})
	}

	return segments
}
