package imageai

type Recommendation struct {
	FormatRecommendations      string `json:"format_recommendations"`
	ResizeRecommendations      string `json:"resize_recommendations"`
	CompressionRecommendations string `json:"compression_recommendations"`
	CachingRecommendations     string `json:"caching_recommendations"`
	AdditionalRecoomendations  string `json:"other_recommendations"`
}
type AIRequest struct {
	Recommendations []Recommendation `json:"rec"`
}
type AIResponse struct {
	Recs []Recommendation `json:"recs"`
}
