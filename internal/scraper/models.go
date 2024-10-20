package scraper

import (
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

func (i Image) String() string {
	return fmt.Sprintf("Image{Src: %s, Alt: %s, Width: %d, Height: %d, Format: %s, Size: %d}", i.Src, i.Alt, i.Width, i.Height, i.Format, i.Size)
}
