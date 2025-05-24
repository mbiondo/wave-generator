package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"
	"time"
	"wave-generator/models"
	"wave-generator/services"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	redisAddr   = getenv("REDIS_ADDR", "localhost:6379")
	redisClient = redis.NewClient(&redis.Options{Addr: redisAddr})
	rateLimit   = 1000 // requests per hour per API key
)

// GenerateAPIKeyHandler generates a new API key (simple UUID, not secure for prod).
func GenerateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	key := "api_" + fmt.Sprintf("%d", time.Now().UnixNano())
	redisClient.Set(context.Background(), "apikey:"+key, "1", 0)
	json.NewEncoder(w).Encode(map[string]string{"api_key": key})
}

// WavePatternHandler processes HTTP requests to extract wave patterns from an image.
// It accepts only POST requests with an image in the request body.
// The function performs the following operations:
// 1. Decodes the image from the request body
// 2. Converts the image to grayscale
// 3. Extracts the wave pattern from the grayscale image
// 4. Fits polynomial segments to represent the pattern
// 5. Generates an SVG representation of the pattern
//
// Returns a JSON response containing:
// - The calculated pattern segments
// - An SVG representation of the pattern
//
// Responds with appropriate HTTP errors if:
// - The request method is not POST (405 Method Not Allowed)
// - The image cannot be decoded (400 Bad Request)
func WavePatternHandler(w http.ResponseWriter, r *http.Request) {
	// Allow same-origin requests without API key
	origin := r.Header.Get("Origin")
	host := r.Host
	isSameOrigin := origin == "" || strings.Contains(origin, host)

	var apiKey string
	if !isSameOrigin {
		apiKey = r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "Missing X-API-Key header", http.StatusUnauthorized)
			return
		}
		ctx := context.Background()
		// Check if API key exists
		exists, err := redisClient.Exists(ctx, "apikey:"+apiKey).Result()
		if err != nil || exists == 0 {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		// Rate limit per API key
		rlKey := "ratelimit:" + apiKey
		count, err := redisClient.Incr(ctx, rlKey).Result()
		if err != nil {
			http.Error(w, "Rate limit error", http.StatusInternalServerError)
			return
		}
		if count == 1 {
			redisClient.Expire(ctx, rlKey, time.Hour)
		}
		if count > int64(rateLimit) {
			ttl, _ := redisClient.TTL(ctx, rlKey).Result()
			http.Error(w, "Rate limit exceeded. Try again in "+ttl.String(), http.StatusTooManyRequests)
			return
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Use POST with image in body", http.StatusMethodNotAllowed)
		return
	}

	img, _, err := image.Decode(r.Body)
	if err != nil {
		http.Error(w, "Error decoding image: "+err.Error(), http.StatusBadRequest)
		return
	}

	var segments []models.PolySegment
	var svg string
	var segmentSVGs []string
	var coords [][]float64

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("processing error: %v", r)
			}
		}()

		b := img.Bounds()
		wImg, hImg := b.Dx(), b.Dy()

		// Convert the image to grayscale and extract the pattern
		gray := services.ToGray(img)
		pattern := services.ExtractPattern(gray, wImg, hImg)
		segments = services.FitSegments(pattern, wImg)
		svg = services.BuildSVG(wImg, hImg, segments)

		// Generate SVG for each segment (mini SVG, width = segment length, height = 40, Y scaled to segment range)
		const miniHeight = 40
		for i := range segments {
			seg := segments[i]
			miniSeg := seg
			miniSeg.X0 = seg.X0
			miniSeg.X1 = seg.X1
			width := seg.X1 - seg.X0 + 1
			if width < 2 {
				segments[i].SVG = ""
				segmentSVGs = append(segmentSVGs, "")
				continue
			}

			// Calcular el rango Y real del segmento en el SVG global
			minY, maxY := seg.CoefA3*float64(seg.X0)*float64(seg.X0)*float64(seg.X0)+seg.CoefA2*float64(seg.X0)*float64(seg.X0)+seg.CoefA1*float64(seg.X0)+seg.CoefA0,
				seg.CoefA3*float64(seg.X1)*float64(seg.X1)*float64(seg.X1)+seg.CoefA2*float64(seg.X1)*float64(seg.X1)+seg.CoefA1*float64(seg.X1)+seg.CoefA0
			for x := seg.X0; x <= seg.X1; x++ {
				y := seg.CoefA3*float64(x)*float64(x)*float64(x) + seg.CoefA2*float64(x)*float64(x) + seg.CoefA1*float64(x) + seg.CoefA0
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
			// Generar SVG solo para el segmento, centrando y escalando Y igual que en el SVG global
			segSVG := services.BuildSVGSegment(seg, width, miniHeight, minY, maxY)
			segmentSVGs = append(segmentSVGs, segSVG)
			segments[i].SVG = segSVG
		}

		// Add coords (pattern as [][x, y])
		coords = make([][]float64, len(pattern))
		for i, y := range pattern {
			coords[i] = []float64{float64(i), y}
		}
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		models.ResponsePayload
		SegmentSVGs []string `json:"segment_svgs"`
	}{
		ResponsePayload: models.ResponsePayload{
			Segments: segments,
			SVG:      svg,
		},
		SegmentSVGs: segmentSVGs,
	})
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
