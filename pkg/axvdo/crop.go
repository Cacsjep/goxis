package axvdo

// CropArea defines the dimensions and position of the crop area.
type CropArea struct {
	Width  int
	Height int
	X      int
	Y      int
}

// calculateCropDimensions calculates the cropping dimensions and position based on input and stream sizes.
func CalculateCropDimensions(inputWidth, inputHeight, streamWidth, streamHeight int) CropArea {
	destWHRatio := float64(inputWidth) / float64(inputHeight)
	cropW := float64(streamWidth)
	cropH := cropW / destWHRatio

	if cropH > float64(streamHeight) {
		cropH = float64(streamHeight)
		cropW = cropH * destWHRatio
	}

	clipW := int(cropW)
	clipH := int(cropH)
	clipX := (streamWidth - clipW) / 2
	clipY := (streamHeight - clipH) / 2

	return CropArea{
		Width:  clipW,
		Height: clipH,
		X:      clipX,
		Y:      clipY,
	}
}
