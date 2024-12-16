package scan

import (
	"net/http"
	"strconv"

	"github.com/voage/sharprender-api/internal/simage"
)

type ScanService struct {
	repo *ScanRepository
}

func NewScanService(repo *ScanRepository) *ScanService {
	return &ScanService{repo: repo}
}

// Parse the URL and store the query params in a FilterOptions struct
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

func applyFilters(images []simage.Image, filters FilterOptions) []simage.Image {
	var filtered []simage.Image

	for _, image := range images {
		// Apply size filter
		if filters.Size != nil && int64(image.Size) <= *filters.Size {
			continue
		}

		// Apply type filter
		if filters.ImgType != nil && image.Format != *filters.ImgType {
			continue
		}

		// Apply load time filter
		if filters.LoadTime != nil && int64(image.Network.LoadTime) <= *filters.LoadTime {
			continue
		}

		// Apply host type filter
		if filters.HostType != nil {
			if *filters.HostType == "first-party" && !isFirstParty(image.Network.InitiatorURL, image.Network.DocumentURL) {
				continue
			}
			if *filters.HostType == "third-party" && isFirstParty(image.Network.InitiatorURL, image.Network.DocumentURL) {
				continue
			}
		}

		filtered = append(filtered, image)
	}

	return filtered
}

func isFirstParty(initiatorURL, documentURL string) bool {
	// Compare the initiator URL and document URL
	return initiatorURL == documentURL
}
