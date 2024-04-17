package axlarod

import "github.com/Cacsjep/goxis/pkg/axvdo"

func CreateCropMap(inputWidth int, inputHeight int, streamWidth int, streamHeight int) (*LarodMap, error) {
	c := axvdo.CalculateCropDimensions(inputWidth, inputHeight, streamWidth, streamHeight)
	return NewLarodMapWithEntries([]*LarodMapEntries{
		{Key: "image.input.crop", Value: [4]int64{int64(c.X), int64(c.Y), int64(c.Width), int64(c.Height)}, ValueType: LarodMapValueTypeIntArr4},
	})
}
