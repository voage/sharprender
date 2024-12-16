package scan

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGetFilterOptions_Print(t *testing.T) {
	// Simulate a query string
	queryParams := "size=50000&type=image/jpeg&loadTime=200&hostType=third-party"

	// Create a mock HTTP request
	req := &http.Request{
		URL: &url.URL{
			RawQuery: queryParams,
		},
	}

	// Call GetFilterOptions
	filters := GetFilterOptions(req)

	// Print the results to the console
	fmt.Println("Parsed FilterOptions:")
	if filters.Size != nil {
		fmt.Printf("Size: %d\n", *filters.Size)
	} else {
		fmt.Println("Size: nil")
	}

	if filters.ImgType != nil {
		fmt.Printf("Type: %s\n", *filters.ImgType)
	} else {
		fmt.Println("Type: nil")
	}

	if filters.LoadTime != nil {
		fmt.Printf("LoadTime: %d\n", *filters.LoadTime)
	} else {
		fmt.Println("LoadTime: nil")
	}

	if filters.HostType != nil {
		fmt.Printf("HostType: %s\n", *filters.HostType)
	} else {
		fmt.Println("HostType: nil")
	}
}
