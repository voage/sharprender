package simage

import (
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
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
	Src     string      `json:"src" bson:"src"`
	Alt     string      `json:"alt" bson:"alt"`
	Width   int         `json:"width" bson:"width"`
	Height  int         `json:"height" bson:"height"`
	Format  string      `json:"format" bson:"format"`
	Size    int         `json:"size" bson:"size"`
	Network NetworkInfo `json:"network" bson:"network"`
}

type NetworkInfo struct {
	RequestID         network.RequestID     `json:"request_id" bson:"request_id"`
	DocumentURL       string                `json:"document_url" bson:"document_url"`
	InitiatorType     network.InitiatorType `json:"initiator_type" bson:"initiator_type"`
	InitiatorURL      string                `json:"initiator_url" bson:"initiator_url"`
	InitiatorLineNo   float64               `json:"initiator_line_no" bson:"initiator_line_no"`
	InitiatorColNo    float64               `json:"initiator_col_no" bson:"initiator_col_no"`
	Method            string                `json:"method" bson:"method"`
	Status            int64                 `json:"status" bson:"status"`
	MimeType          string                `json:"mime_type" bson:"mime_type"`
	Protocol          string                `json:"protocol" bson:"protocol"`
	RemoteIPAddress   string                `json:"remote_ip_address" bson:"remote_ip_address"`
	RemotePort        int64                 `json:"remote_port" bson:"remote_port"`
	EncodedDataLength int                   `json:"encoded_data_length" bson:"encoded_data_length"`
	RequestTime       *cdp.MonotonicTime    `json:"request_time" bson:"request_time"`
	ResponseTime      *cdp.MonotonicTime    `json:"response_time" bson:"response_time"`
	LoadTime          float64               `json:"load_time" bson:"load_time"`
}

type ImageOverview struct {
	TotalImages         int            `json:"total_images" bson:"total_images"`
	TotalSize           int            `json:"total_size" bson:"total_size"`
	AverageSize         int            `json:"average_size" bson:"average_size"`
	AverageWidth        int            `json:"average_width" bson:"average_width"`
	AverageHeight       int            `json:"average_height" bson:"average_height"`
	Formats             map[string]int `json:"formats" bson:"formats"`
	CacheHits           int            `json:"cache_hits" bson:"cache_hits"`
	AverageRequestTime  float64        `json:"average_request_time" bson:"average_request_time"`
	AverageResponseTime float64        `json:"average_response_time" bson:"average_response_time"`
	AverageTotalTime    float64        `json:"average_total_time" bson:"average_total_time"`
}

type NetworkProfile struct {
	Download float64
	Upload   float64
	Latency  float64
}
