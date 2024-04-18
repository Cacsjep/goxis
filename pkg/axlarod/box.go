package axlarod

// BoundingBox is a struct that holds the coordinates of a bounding box in float32
type BoundingBox struct {
	Top    float32
	Left   float32
	Bottom float32
	Right  float32
}

// Cords32 is a struct that holds the coordinates of a bounding box in float32
type Cords32 struct {
	X float32
	Y float32
	W float32
	H float32
}

// Cords64 is a struct that holds the coordinates of a bounding box in float64
type Cords64 struct {
	X float64
	Y float64
	W float64
	H float64
}

// Scale scales the bounding box to the given width and height
func (bbox *BoundingBox) Scale(w, h int) BoundingBox {
	return BoundingBox{
		Top:    bbox.Top * float32(h),
		Left:   bbox.Left*float32(h) + float32(w-h)/2,
		Bottom: bbox.Top*float32(h) + (bbox.Bottom-bbox.Top)*float32(h),
		Right:  bbox.Left*float32(h) + (bbox.Right-bbox.Left)*float32(h) + float32(w-h)/2,
	}
}

// ToCords32 converts the bounding box to Cords32
func (bbox *BoundingBox) ToCords32() Cords32 {
	return Cords32{X: bbox.Left, Y: bbox.Top, W: bbox.Right - bbox.Left, H: bbox.Bottom - bbox.Top}
}

// ToCords64 converts the bounding box to Cords64
func (bbox *BoundingBox) ToCords64() Cords64 {
	return Cords64{X: float64(bbox.Left), Y: float64(bbox.Top), W: float64(bbox.Right - bbox.Left), H: float64(bbox.Bottom - bbox.Top)}
}
