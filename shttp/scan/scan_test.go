package scan

import (
	"testing"

	"github.com/voage/sharprender-api/internal/simage"
)

var mockImages = []simage.Image{
	{
		Src:    "https://example.com/image1.jpg",
		Format: "image/jpeg",
		Size:   10000,
		Network: simage.NetworkInfo{
			InitiatorURL: "https://example.com",
			DocumentURL:  "https://example.com",
			LoadTime:     150,
		},
	},
	{
		Src:    "https://example.com/image2.png",
		Format: "image/png",
		Size:   50000,
		Network: simage.NetworkInfo{
			InitiatorURL: "https://thirdparty.com",
			DocumentURL:  "https://example.com",
			LoadTime:     300,
		},
	},
	{
		Src:    "https://example.com/image3.gif",
		Format: "image/gif",
		Size:   20000,
		Network: simage.NetworkInfo{
			InitiatorURL: "https://example.com",
			DocumentURL:  "https://example.com",
			LoadTime:     50,
		},
	},
}

func TestApplyFilters_SizeFilter(t *testing.T) {
	filters := FilterOptions{
		Size: &[]int64{30000}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 1 {
		t.Errorf("Expected 1 image, got %d", len(result))
	}

	if result[0].Src != "https://example.com/image2.png" {
		t.Errorf("Expected 'image2.png', got '%s'", result[0].Src)
	}
}

func TestApplyFilters_TypeFilter(t *testing.T) {
	filters := FilterOptions{
		ImgType: &[]string{"image/jpeg"}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 1 {
		t.Errorf("Expected 1 image, got %d", len(result))
	}

	if result[0].Src != "https://example.com/image1.jpg" {
		t.Errorf("Expected 'image1.jpg', got '%s'", result[0].Src)
	}
}

func TestApplyFilters_LoadTimeFilter(t *testing.T) {
	filters := FilterOptions{
		LoadTime: &[]int64{200}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 1 {
		t.Errorf("Expected 2 images, got %d", len(result))
	}

	if result[0].Src != "https://example.com/image2.png" {
		t.Errorf("Expected 'image2.jpg', got '%s'", result[0].Src)
	}

}

func TestApplyFilters_HostTypeFilter_FirstParty(t *testing.T) {
	filters := FilterOptions{
		HostType: &[]string{"first-party"}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 images, got %d", len(result))
	}

	if result[0].Src != "https://example.com/image1.jpg" {
		t.Errorf("Expected 'image1.jpg', got '%s'", result[0].Src)
	}

	if result[1].Src != "https://example.com/image3.gif" {
		t.Errorf("Expected 'image3.gif', got '%s'", result[1].Src)
	}
}

func TestApplyFilters_HostTypeFilter_ThirdParty(t *testing.T) {
	filters := FilterOptions{
		HostType: &[]string{"third-party"}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 1 {
		t.Errorf("Expected 1 image, got %d", len(result))
	}

	if result[0].Src != "https://example.com/image2.png" {
		t.Errorf("Expected 'image2.png', got '%s'", result[0].Src)
	}
}

func TestApplyFilters_MultipleFilters(t *testing.T) {
	filters := FilterOptions{
		Size:     &[]int64{15000}[0],
		ImgType:  &[]string{"image/jpeg"}[0],
		HostType: &[]string{"first-party"}[0],
	}

	result := applyFilters(mockImages, filters)

	// Check the number of images returned
	if len(result) != 0 {
		t.Errorf("Expected 0 images, got %d", len(result))
	}
}

func TestApplyFilters_NoFilters(t *testing.T) {
	filters := FilterOptions{}

	result := applyFilters(mockImages, filters)

	if len(result) != len(mockImages) {
		t.Errorf("Expected %d images, got %d", len(mockImages), len(result))
	}
}

func TestApplyFilters_NoMatches(t *testing.T) {
	filters := FilterOptions{
		Size: &[]int64{1000000}[0],
	}

	result := applyFilters(mockImages, filters)

	if len(result) != 0 {
		t.Errorf("Expected 0 images, got %d", len(result))
	}
}
