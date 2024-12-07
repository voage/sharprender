package scan

import (
	"time"

	"github.com/voage/sharprender-api/internal/simage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scan struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	URL       string             `json:"url" bson:"url"`
	Images    []simage.Image     `json:"images" bson:"images"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
