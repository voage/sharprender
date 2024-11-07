package simage

import (
	"os"
	"testing"
)

func TestCompressImages(t *testing.T) {
	image := Image{
		Src: "https://www.superherotoystore.com/cdn/shop/articles/e33c2fa94c03efa06678116f80d62d0d_708x.jpg?v=1590599656",
		Alt: "test-image",
	}

	params := ImageParams{
		Quality: 80,
		Width:   800,
		Height:  600,
	}

	err := CompressImages(params, image)
	if err != nil {
		t.Fatalf("CompressImages failed: %v", err)
	}

	outputPath := "compressed_test-image.webp"

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Expected output file %s not found", outputPath)
	} else {
		t.Logf("Output file %s created successfully", outputPath)
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Logf("Warning: failed to remove test file %s", outputPath)
	}
}
