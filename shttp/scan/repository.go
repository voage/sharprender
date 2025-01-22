package scan

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *ScanRepository) Create(ctx context.Context, scan *Scan) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, scan)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// FindWithFilter fetches a scan document and filters its images using aggregation
func (r *ScanRepository) FindWithFilter(ctx context.Context, scanFilter, imageFilter bson.M) (*Scan, error) {
	// MongoDB aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: scanFilter}},  // Match scan document by ID
		{{Key: "$unwind", Value: "$images"}},  // Unwind the images array
		{{Key: "$match", Value: imageFilter}}, // Apply image-level filters
		{
			{Key: "$group", Value: bson.M{ // Regroup results into a single document
				"_id":       "$_id",
				"URL":       bson.M{"$first": "$URL"},
				"images":    bson.M{"$push": "$images"},
				"createdAt": bson.M{"$first": "$createdAt"},
			}},
		},
	}

	// Execute aggregation
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	// Decode the result into a Scan struct
	var results []Scan
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Return the first (and only) result or a "not found" error
	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}
