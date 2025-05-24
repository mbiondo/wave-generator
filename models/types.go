package models

type PolySegment struct {
	X0         int     `json:"domain_start"`
	X1         int     `json:"domain_end"`
	CoefA3     float64 `json:"a3"`
	CoefA2     float64 `json:"a2"`
	CoefA1     float64 `json:"a1"`
	CoefA0     float64 `json:"a0"`
	Expression string  `json:"expression"`
	SVG        string  `json:"svg,omitempty"`
}

type ResponsePayload struct {
	Segments []PolySegment `json:"segments"`
	SVG      string        `json:"svg"`
}
