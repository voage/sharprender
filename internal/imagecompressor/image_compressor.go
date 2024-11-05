package imagecompressor

import (
	"fmt"
	"io"
	"net/http"

	"github.com/h2non/bimg"
	"github.com/voage/sharprender-api/internal/imagescraper"
)

func CompressImages(ip ImageParams, i imagescraper.Image) error {

	resp, err := http.Get(i.Src)
	if err != nil {
		return fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching image %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}

	options := bimg.Options{
		Width:   ip.Width,
		Height:  ip.Height,
		Quality: ip.Quality,
		Type:    bimg.WEBP,
	}

	newImage, err := bimg.NewImage(imageData).Process(options)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	outputPath := "compressed_" + i.Alt + ".webp" 
	err = bimg.Write(outputPath, newImage)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}

	fmt.Printf("Image compressed and saved to %s\n", outputPath)
	return nil
}
