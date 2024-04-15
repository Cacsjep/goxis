package axlarod

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

// ConvertRGBToImage converts raw RGB bytes to an image.Image object
func ConvertRGBToImage(rgb []byte, width, height int) (*image.RGBA, error) {
	if len(rgb) != width*height*3 {
		return nil, fmt.Errorf("invalid RGB data length: got %d, expected %d", len(rgb), width*height*3)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := rgb[index]
			g := rgb[index+1]
			b := rgb[index+2]
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			index += 3
		}
	}

	return img, nil
}

// SaveImageAsJPEG saves an image.Image as a JPEG file
func SaveImageAsJPEG(rgb []byte, width, height int, filename string) error {
	img, err := ConvertRGBToImage(rgb, width, height)
	if err != nil {
		return fmt.Errorf("failed to convert RGB to image: %v", err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return fmt.Errorf("failed to encode image to JPEG: %v", err)
	}

	return nil
}
