package imagescraper

func GetImageOverview(images []Image) ImageOverview {
	overview := ImageOverview{
		TotalImages: len(images),
		Formats:     make(map[string]int),
	}

	var totalWidth, totalHeight, totalSize, totalCacheHits int
	var totalRequestTime, totalResponseTime, totalTiming float64

	for _, img := range images {
		totalSize += img.Size
		overview.Formats[img.Format]++

		if img.Cache.CacheHit {
			totalCacheHits++
		}

		totalWidth += img.Width
		totalHeight += img.Height

		if img.Network.RequestTime != nil && img.Network.ResponseTime != nil {
			requestTime := float64(img.Network.RequestTime.Time().Unix())
			responseTime := float64(img.Network.ResponseTime.Time().Unix())
			totalRequestTime += responseTime - requestTime
			totalResponseTime += responseTime
		}

		totalTiming += img.Timing.TotalTime
	}

	if overview.TotalImages > 0 {
		overview.AverageSize = totalSize / overview.TotalImages
		overview.AverageWidth = totalWidth / overview.TotalImages
		overview.AverageHeight = totalHeight / overview.TotalImages
		overview.TotalSize = totalSize
		overview.CacheHits = totalCacheHits

		overview.AverageRequestTime = totalRequestTime / float64(overview.TotalImages)
		overview.AverageResponseTime = totalResponseTime / float64(overview.TotalImages)
		overview.AverageTotalTime = totalTiming / float64(overview.TotalImages)
	}

	return overview
}
