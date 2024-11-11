package scan

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ScanRepository struct {
	collection *mongo.Collection
}

func NewScanRepository(client *mongo.Client) *ScanRepository {
	return &ScanRepository{
		collection: client.Database("sharprenderdb").Collection("scans"),
	}
}

func (r *ScanRepository) FindOne(ctx context.Context, filter interface{}) (*Scan, error) {
	var scan Scan
	err := r.collection.FindOne(ctx, filter).Decode(&scan)
	return &scan, err
}

func (r *ScanRepository) Create(ctx context.Context, scan *Scan) error {
	_, err := r.collection.InsertOne(ctx, scan)
	return err
}
