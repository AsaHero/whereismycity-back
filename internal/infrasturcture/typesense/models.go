package typesense

type Locations struct {
	ID      int64   `json:"id"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	Country string  `json:"country"`
	Code    string  `json:"code"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`

	VectorDistance  *float32 `json:"_vector_distance"`
	TextMatchScore  *int64   `json:"text_match"`
	RankFusionScore *float64 `json:"_rank_fusion_score"`
}
