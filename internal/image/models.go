package image

import (
	"github.com/chromedp/cdproto/cdp"
)

type Recommendation struct {
	FormatRecommendations      string `json:"format_recommendations"`
	ResizeRecommendations      string `json:"resize_recommendations"`
	CompressionRecommendations string `json:"compression_recommendations"`
	CachingRecommendations     string `json:"caching_recommendations"`
	AdditionalRecommendations  string `json:"other_recommendations"`
}
type AIRequest struct {
	Recommendations []Recommendation `json:"rec"`
}
type AIResponse struct {
	Recs []Recommendation `json:"recs"`
}
type ImageParams struct {
	Quality int
	Width   int
	Height  int
}

type Image struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Size   int    `json:"size"`

	Network NetworkInfo `json:"network"`
	Timing  TimingInfo  `json:"timing"`
	Cache   CacheInfo   `json:"cache"`
}

type NetworkInfo struct {
	RequestID         string              `json:"request_id"`
	DocumentURL       string              `json:"document_url"`
	InitiatorType     string              `json:"initiator_type"`
	InitiatorURL      string              `json:"initiator_url"`
	InitiatorLineNo   int                 `json:"initiator_line_no"`
	InitiatorColNo    int                 `json:"initiator_col_no"`
	Method            string              `json:"method"`
	Status            int                 `json:"status"`
	MimeType          string              `json:"mime_type"`
	Protocol          string              `json:"protocol"`
	RemoteIPAddress   string              `json:"remote_ip_address"`
	RemotePort        int                 `json:"remote_port"`
	EncodedDataLength int                 `json:"encoded_data_length"`
	RequestTime       *cdp.TimeSinceEpoch `json:"request_time"`
	ResponseTime      *cdp.MonotonicTime  `json:"response_time"`
}

type TimingInfo struct {
	DNSTime     float64 `json:"dns_time"`
	ConnectTime float64 `json:"connect_time"`
	SSLTime     float64 `json:"ssl_time"`
	SendTime    float64 `json:"send_time"`
	WaitTime    float64 `json:"wait_time"`
	ReceiveTime float64 `json:"receive_time"`
	TotalTime   float64 `json:"total_time"`
}

type CacheInfo struct {
	FromCache       bool   `json:"from_cache"`
	CacheHit        bool   `json:"cache_hit"`
	CacheState      string `json:"cache_state,omitempty"`
	CacheValidation string `json:"cache_validation,omitempty"`
	Age             int    `json:"age,omitempty"`
	ExpirationTime  string `json:"expiration_time,omitempty"`
	LastModified    string `json:"last_modified,omitempty"`
	ETag            string `json:"etag,omitempty"`
}

type ImageOverview struct {
	TotalImages         int            `json:"total_images"`
	TotalSize           int            `json:"total_size"`
	AverageSize         int            `json:"average_size"`
	AverageWidth        int            `json:"average_width"`
	AverageHeight       int            `json:"average_height"`
	Formats             map[string]int `json:"formats"`
	CacheHits           int            `json:"cache_hits"`
	AverageRequestTime  float64        `json:"average_request_time"`
	AverageResponseTime float64        `json:"average_response_time"`
	AverageTotalTime    float64        `json:"average_total_time"`
}
