package scraper

import (
	"github.com/chromedp/cdproto/cdp"
)

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
