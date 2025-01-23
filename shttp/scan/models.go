package scan

import (
	"time"

	"github.com/voage/sharprender-api/internal/simage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ScanID    string             `json:"scan_id" bson:"scan_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	URL       string             `json:"url" bson:"url"`
	Images    []simage.Image     `json:"images" bson:"images"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FilterOptions struct {
	Size     *int64
	ImgType  *string
	LoadTime *int64
	HostType *string
}

type ScanResult struct {
	Images       []simage.Image         `json:"images"`
	Aggregations map[string]interface{} `json:"aggregations"`
}
