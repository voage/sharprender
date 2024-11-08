package simage

import (
	"fmt"
	"io"
	"net/http"

	"github.com/h2non/bimg"
)

func CompressImages(ip ImageParams, i Image) error {

	options := bimg.Options{
		Width:   ip.Width,
		Height:  ip.Height,
		Quality: ip.Quality,
		Type:    bimg.WEBP,
	}

	imageData, err := fetchImageData(i.Src)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}

	newImage, err := bimg.NewImage(imageData).Process(options)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	err = saveImages(newImage, i.Alt)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	return nil
}

func fetchImageData(src string) ([]byte, error) {
	resp, err := http.Get(src)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching image %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageData, nil
}

func saveImages(imageData []byte, alt string) error {
	outputPath := "compressed_" + alt + ".webp"
	err := bimg.Write(outputPath, imageData)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}
	fmt.Printf("Image compressed and saved to %s\n", outputPath)
	return nil
}
