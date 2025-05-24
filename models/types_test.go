package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPolySegmentJSONMarshalling(t *testing.T) {
	segment := PolySegment{
		X0:         0,
		X1:         10,
		CoefA3:     1.0,
		CoefA2:     2.0,
		CoefA1:     3.0,
		CoefA0:     4.0,
		Expression: "x^3 + 2x^2 + 3x + 4",
		SVG:        "<svg></svg>",
	}

	data, err := json.Marshal(segment)
	if err != nil {
		t.Fatalf("Failed to marshal PolySegment: %v", err)
	}

	var unmarshalled PolySegment
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal PolySegment: %v", err)
	}

	if !reflect.DeepEqual(segment, unmarshalled) {
		t.Errorf("Expected %v, got %v", segment, unmarshalled)
	}
}

func TestResponsePayloadJSONMarshalling(t *testing.T) {
	payload := ResponsePayload{
		Segments: []PolySegment{
			{
				X0:         0,
				X1:         5,
				CoefA3:     0.5,
				CoefA2:     1.5,
				CoefA1:     2.5,
				CoefA0:     3.5,
				Expression: "0.5x^3 + 1.5x^2 + 2.5x + 3.5",
				SVG:        "",
			},
		},
		SVG: "<svg>payload</svg>",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ResponsePayload: %v", err)
	}

	var unmarshalled ResponsePayload
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal ResponsePayload: %v", err)
	}

	if !reflect.DeepEqual(payload, unmarshalled) {
		t.Errorf("Expected %v, got %v", payload, unmarshalled)
	}
}
