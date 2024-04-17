package axlarod

type BoundingBox struct {
	Top    float32
	Left   float32
	Bottom float32
	Right  float32
}

type Cords32 struct {
	X float32
	Y float32
	W float32
	H float32
}

type Cords64 struct {
	X float64
	Y float64
	W float64
	H float64
}

func (bbox *BoundingBox) Scale(w, h int) BoundingBox {
	return BoundingBox{
		Top:    bbox.Top * float32(h),
		Left:   bbox.Left*float32(h) + float32(w-h)/2,                                     // This simulates the (widthFrameHD - heightFrameHD) / 2 offset
		Bottom: bbox.Top*float32(h) + (bbox.Bottom-bbox.Top)*float32(h),                   // Bottom now represents absolute position
		Right:  bbox.Left*float32(h) + (bbox.Right-bbox.Left)*float32(h) + float32(w-h)/2, // Right now represents absolute position
	}
}

func (bbox *BoundingBox) ToCords32() Cords32 {
	return Cords32{X: bbox.Left, Y: bbox.Top, W: bbox.Right - bbox.Left, H: bbox.Bottom - bbox.Top}
}

func (bbox *BoundingBox) ToCords64() Cords64 {
	return Cords64{X: float64(bbox.Left), Y: float64(bbox.Top), W: float64(bbox.Right - bbox.Left), H: float64(bbox.Bottom - bbox.Top)}
}
