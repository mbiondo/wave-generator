package services

import (
	"strings"
	"testing"
	"wave-generator/models"
)

func TestBuildSVG(t *testing.T) {
	segments := []models.PolySegment{
		{X0: 0, X1: 4, CoefA3: 0, CoefA2: 0, CoefA1: 1, CoefA0: 0},
		{X0: 5, X1: 9, CoefA3: 0, CoefA2: 0, CoefA1: 1, CoefA0: 5},
	}

	svg := BuildSVG(10, 10, segments)

	if svg == "" {
		t.Fatalf("expected non-empty SVG, got empty string")
	}

	if !contains(svg, "<svg") || !contains(svg, "</svg>") {
		t.Errorf("invalid SVG format: %s", svg)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
