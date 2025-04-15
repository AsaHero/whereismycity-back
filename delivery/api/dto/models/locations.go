package models

type Location struct {
	ID        int64   `json:"id"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	VectorDistance  *float32 `json:"vector_distance"`
	TextMatchScore  *int64   `json:"text_match_score"`
	RankFusionScore *float64 `json:"rank_fusion_score"`
}

type SearchRequest struct {
	Query string `form:"query" validate:"required,min=3,max=100"`
	Limit uint   `form:"limit" validate:"required,min=1,max=100,default=10"`
}

type SearchResponse struct {
	Query     string      `json:"query"`
	Limit     uint        `json:"limit"`
	Locations []*Location `json:"locations"`
}
