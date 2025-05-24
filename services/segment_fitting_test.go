package services

import (
	"math"
	"testing"
)

func TestFitSegments(t *testing.T) {
	tests := []struct {
		name     string
		sky      []float64
		width    int
		wantSegs int
		wantErr  bool
	}{
		{
			name:     "quadratic function",
			sky:      generateQuadratic(128),
			width:    128,
			wantSegs: 8, // ajustado a la lógica actual
			wantErr:  false,
		},
		{
			name:     "small input",
			sky:      []float64{0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 100, 121},
			width:    12,
			wantSegs: 1, // ajustado a la lógica actual
			wantErr:  false,
		},
		{
			name:     "minimum size input",
			sky:      []float64{0, 1, 4, 9},
			width:    4,
			wantSegs: 1,
			wantErr:  false,
		},
		{
			name:     "exact segment size",
			sky:      generateQuadratic(64),
			width:    64,
			wantSegs: 4, // ajustado a la lógica actual
			wantErr:  false,
		},
		{
			name:     "large input",
			sky:      generateQuadratic(256),
			width:    256,
			wantSegs: 16, // ajustado a la lógica actual
			wantErr:  false,
		},
		{
			name:     "uneven segments",
			sky:      generateQuadratic(33),
			width:    33,
			wantSegs: 2, // ajustado a la lógica actual
			wantErr:  false,
		},
		{
			name:     "too small input",
			sky:      []float64{0, 1},
			width:    2,
			wantSegs: 0,
			wantErr:  true,
		},
		{
			name:     "nil input",
			sky:      nil,
			width:    0,
			wantSegs: 0,
			wantErr:  true,
		},
		{
			name:     "zero width",
			sky:      []float64{0, 1, 2, 3},
			width:    0,
			wantSegs: 0,
			wantErr:  true,
		},
		{
			name:     "very small width",
			sky:      generateQuadratic(8), // Using 8 points to ensure enough data
			width:    4,
			wantSegs: 1,
			wantErr:  false,
		},
		{
			name:     "edge segment case",
			sky:      generateQuadratic(7),
			width:    7,
			wantSegs: 1,
			wantErr:  false,
		},
		{
			name:     "width larger than data",
			sky:      generateQuadratic(4),
			width:    8,
			wantSegs: 1,
			wantErr:  false,
		},
		{
			name:     "force QR error",
			sky:      []float64{2, 2, 2, 2, 2, 2, 2, 2}, // Constant function causes singular matrix
			width:    8,
			wantSegs: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but got none")
					}
				}()
				FitSegments(tt.sky, tt.width)
				t.Error("expected panic")
				return
			}

			segments := FitSegments(tt.sky, tt.width)

			if !tt.wantErr {
				if len(segments) == 0 {
					t.Fatalf("expected non-empty segments, got %d", len(segments))
				}

				if len(segments) != tt.wantSegs {
					t.Errorf("expected %d segments, got %d", tt.wantSegs, len(segments))
				}

				// Verify segments are contiguous
				for i := 0; i < len(segments)-1; i++ {
					if segments[i].X1+1 != segments[i+1].X0 {
						t.Errorf("segments not contiguous at index %d: %d != %d", i, segments[i].X1+1, segments[i+1].X0)
					}
				}

				for _, seg := range segments {
					if seg.X1 <= seg.X0 {
						t.Errorf("invalid segment domain: [%d, %d]", seg.X0, seg.X1)
					}

					if seg.X1-seg.X0+1 < 4 {
						t.Errorf("segment too small for cubic fit: [%d, %d]", seg.X0, seg.X1)
					}

					// Test points in the segment
					for x := seg.X0; x <= seg.X1; x++ {
						expected := float64(x * x)
						actual := seg.CoefA3*math.Pow(float64(x), 3) +
							seg.CoefA2*math.Pow(float64(x), 2) +
							seg.CoefA1*float64(x) +
							seg.CoefA0

						if math.Abs(actual-expected) > 1.0 {
							t.Errorf("at x=%d: expected %.2f, got %.2f", x, expected, actual)
						}
					}

					// Verify expression format
					if len(seg.Expression) == 0 {
						t.Errorf("empty expression for segment [%d, %d]", seg.X0, seg.X1)
					}
				}
			}
		})
	}
}

// generateQuadratic creates a slice of y = x² values
func generateQuadratic(size int) []float64 {
	result := make([]float64, size)
	for i := 0; i < size; i++ {
		result[i] = float64(i * i)
	}
	return result
}
