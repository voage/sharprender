package scan

import (
	"net/http"
	"strconv"
)

type ScanService struct {
	repo *ScanRepository
}

func NewScanService(repo *ScanRepository) *ScanService {
	return &ScanService{repo: repo}
}

func GetFilterOptions(r *http.Request) FilterOptions {

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
