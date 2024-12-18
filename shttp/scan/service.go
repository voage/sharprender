package scan

import (
	"context"
	"net/http"
	"strconv"

	"github.com/voage/sharprender-api/internal/simage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScanService struct {
	repo *ScanRepository
}

func NewScanService(repo *ScanRepository) *ScanService {
	return &ScanService{repo: repo}
}

// Parse the URL and store the query params in a FilterOptions struct
func parseFilterOptions(r *http.Request) FilterOptions {


	var filters FilterOptions

	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return filters
		}
		filters.Size = &size
	}

	if imageType := r.URL.Query().Get("type"); imageType != "" {
		filters.ImgType = &imageType
	}

	if loadTimeStr := r.URL.Query().Get("loadTime"); loadTimeStr != "" {
		loadTime, err := strconv.ParseInt(loadTimeStr, 10, 64)
		if err != nil {
			return filters
		}
		filters.LoadTime = &loadTime
	}

	if hostType := r.URL.Query().Get("hostType"); hostType != "" {
		filters.HostType = &hostType
	}

	return filters
}

// buildMongoFilterQuery builds a MongoDB query for filtering images
func buildMongoFilterQuery(filters FilterOptions) bson.M {
	query := bson.M{}

	if filters.Size != nil {
		query["images.size"] = bson.M{"$gt": *filters.Size}
	}
	if filters.ImgType != nil {
		query["images.format"] = *filters.ImgType
	}
	if filters.LoadTime != nil {
		query["images.network.loadTime"] = bson.M{"$gt": *filters.LoadTime}
	}
	if filters.HostType != nil {
		if *filters.HostType == "first-party" {
			query["images.network.initiatorURL"] = bson.M{"$eq": "images.network.documentURL"}
		} else if *filters.HostType == "third-party" {
			query["images.network.initiatorURL"] = bson.M{"$ne": "images.network.documentURL"}
		}
	}

	return query
}

// calculateAggregations calculates metrics for all images
func calculateAggregations(images []simage.Image) map[string]interface{} {
	var totalSize, totalLoadTime int64
	var avgSize, avgLoadTime float64

	for _, img := range images {
		totalSize += int64(img.Size)
		totalLoadTime += int64(img.Network.LoadTime)
	}

	count := len(images)
	if count > 0 {
		avgSize = float64(totalSize) / float64(count)
		avgLoadTime = float64(totalLoadTime) / float64(count)
	}

	return map[string]interface{}{
		"avgSize":       avgSize,
		"totalSize":     totalSize,
		"avgLoadTime":   avgLoadTime,
		"totalLoadTime": totalLoadTime,
		"imageCount":    count,
	}
}

func (s *ScanService) fetchScanResult(ctx context.Context, id primitive.ObjectID, filters FilterOptions) (*ScanResult, error) {
	imageFilter := buildMongoFilterQuery(filters)

	// Fetch scan with filtered images
	scan, err := s.repo.FindWithFilter(ctx, bson.M{"_id": id}, imageFilter)
	if err != nil {
		return nil, err
	}

	// Calculate aggregations for all images
	aggregations := calculateAggregations(scan.Images)

	return &ScanResult{
		Images:       scan.Images,
		Aggregations: aggregations,
	}, nil
}
