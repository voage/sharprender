package imageai

type Recommendation struct {
	Title string `json:"title"`
	Details string `json:"details"`

}
type AIRequest struct {
	Recommendations []Recommendation `json:"rec"`
}
type AIResponse struct {
	Recs []Recommendation `json:"recs"`
}
