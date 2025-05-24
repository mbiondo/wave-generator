package services

import (
	"fmt"
	"wave-generator/models"
)

// BuildSVG generates an SVG representation of a polynomial curve defined by segments.
//
// It creates an SVG image with the specified width and height, containing a polyline
// that represents the polynomial function described by the provided segments.
// The function evaluates each point along the x-axis, finds the appropriate polynomial segment,
// and calculates the corresponding y-value using the segment's coefficients.
//
// Parameters:
//   - w: Width of the SVG image in pixels
//   - h: Height of the SVG image in pixels
//   - segs: Slice of polynomial segments defining the curve
//
// Returns:
//   - A string containing the complete SVG markup with the plotted polyline
//
// Each segment defines a cubic polynomial of the form:
//
//	y = a3*x^3 + a2*x^2 + a1*x + a0
//
// where the coefficients (a3, a2, a1, a0) are provided in each PolySegment struct.
func BuildSVG(w, h int, segs []models.PolySegment) string {
	s := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg"><polyline fill="none" stroke="lime" stroke-width="1" points="`, w, h)
	for x := 0; x < w; x++ {
		var seg models.PolySegment
		for _, s := range segs {
			if x >= s.X0 && x <= s.X1 {
				seg = s
				break
			}
		}
		xf := float64(x)
		// Expand math.Pow for linter
		y := seg.CoefA3*xf*xf*xf + seg.CoefA2*xf*xf + seg.CoefA1*xf + seg.CoefA0
		s += fmt.Sprintf("%d,%.2f ", x, y)
	}
	s += `"/></svg>`
	return s
}

// BuildSVGSegment generates an SVG for a single segment, scaling Y to [minY, maxY] and X to [0,width]
func BuildSVGSegment(seg models.PolySegment, width, height int, minY, maxY float64) string {
	points := ""
	yRange := maxY - minY
	if yRange == 0 {
		yRange = 1
	}
	for i := 0; i < width; i++ {
		x := seg.X0 + i
		y := seg.CoefA3*float64(x)*float64(x)*float64(x) + seg.CoefA2*float64(x)*float64(x) + seg.CoefA1*float64(x) + seg.CoefA0
		// Rotar en X e Y: invertir X y Y respecto al centro
		px := width - 1 - i
		py := float64(height-2) - ((y-minY)/yRange)*float64(height-4)
		py = float64(height-2) - py // invertir Y
		points += fmt.Sprintf("%d,%.1f ", px, py)
	}
	return fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" style="background:#f8fafd;border-radius:4px;border:1px solid #e1e4e8;"><polyline fill="none" stroke="#3498db" stroke-width="2" points="%s"/></svg>`, width, height, width, height, points)
}
