package scraper

import (
	"errors"
	"fmt"
)

type Image struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

func NewImage(src, alt string, width, height int, format string, size int) (*Image, error) {
	if src == "" {
		return nil, errors.New("src cannot be empty")
	}

	return &Image{
		Src:    src,
		Alt:    alt,
		Width:  width,
		Height: height,
		Format: format,
		Size:   size,
	}, nil
}

func (i *Image) String() string {
	return fmt.Sprintf("Image{Src: %s, Alt: %s, Width: %d, Height: %d, Format: %s, Size: %d}", i.Src, i.Alt, i.Width, i.Height, i.Format, i.Size)
}
