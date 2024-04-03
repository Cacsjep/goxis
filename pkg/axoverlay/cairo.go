package axoverlay

/*
#cgo LDFLAGS: -laxoverlay
#cgo pkg-config: gio-2.0 glib-2.0 cairo axoverlay
#include <axoverlay.h>
#include <cairo/cairo.h>
*/
import "C"
import (
	"image/color"
	"unsafe"
)

const PALETTE_VALUE_RANGE = 255.0

type CairoContext struct {
	ptr *C.cairo_t
}

func NewCairoCtxFromC(renderingContext C.gpointer) *CairoContext {
	return &CairoContext{ptr: (*C.cairo_t)(renderingContext)}
}

func Index2Cairo(colorIndex int) float64 {
	return (float64((colorIndex << 4) + colorIndex)) / PALETTE_VALUE_RANGE
}

func (ctx *CairoContext) DrawArrow(x, y, length, angle float64, color color.RGBA) {
	C.cairo_save(ctx.ptr)
	ctx.Translate(x, y)
	ctx.Rotate(angle * (3.141592653589793 / 180.0))
	ctx.SetSourceRGB(color)
	ctx.SetLineWidth(6)
	ctx.MoveTo(0, 0)
	shaft := float64(length / 4)
	arrowhead := float64(9)
	ctx.LineTo(length-shaft, 0)       // Draw the line for the arrow's shaft
	ctx.RelLineTo(-shaft, -arrowhead) // Left part of the arrowhead
	ctx.RelMoveTo(shaft, arrowhead)   // Move back to the tip of the arrow
	ctx.RelLineTo(-shaft, arrowhead)  // Right part of the arrowhead
	ctx.Stroke()
	C.cairo_restore(ctx.ptr)
}

func (ctx *CairoContext) DrawText(text string, x float64, y float64, size float64, font_name string, color color.RGBA) {
	ctx.SetSourceRGB(color)
	ctx.SelectFontFace(font_name, FONT_SLANT_NORMAL, FONT_WEIGHT_NORMAL)
	ctx.SetFontSize(size)
	e := ctx.TextExtents(text)
	ctx.MoveTo(x, y+e.Height)
	ctx.ShowText(text)
}

func (ctx *CairoContext) DrawRect(x float64, y float64, width float64, height float64, color color.RGBA, linewidth float64) {
	ctx.SetSourceRGBA(color)
	ctx.SetOperator(OPERATOR_SOURCE)
	ctx.SetLineWidth(linewidth)
	ctx.Rectangle(x, y, width, height)
	ctx.Stroke()
}

func (ctx *CairoContext) DrawTransparent(width, height int) {
	ctx.SetSourceRGBA(ColorTransparent)
	ctx.SetOperator(OPERATOR_SOURCE)
	ctx.Rectangle(0, 0, float64(width), float64(height))
	ctx.Fill()
}

func (ctx *CairoContext) NewPath() {
	C.cairo_new_path(ctx.ptr)
}

func (ctx *CairoContext) MoveTo(x, y float64) {
	C.cairo_move_to(ctx.ptr, C.double(x), C.double(y))
}

func (ctx *CairoContext) NewSubPath() {
	C.cairo_new_sub_path(ctx.ptr)
}

func (ctx *CairoContext) LineTo(x, y float64) {
	C.cairo_line_to(ctx.ptr, C.double(x), C.double(y))
}

func (ctx *CairoContext) CurveTo(x1, y1, x2, y2, x3, y3 float64) {
	C.cairo_curve_to(ctx.ptr,
		C.double(x1), C.double(y1),
		C.double(x2), C.double(y2),
		C.double(x3), C.double(y3))
}

func (ctx *CairoContext) Arc(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc(ctx.ptr,
		C.double(xc), C.double(yc),
		C.double(radius),
		C.double(angle1), C.double(angle2))
}

func (ctx *CairoContext) ArcNegative(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc_negative(ctx.ptr,
		C.double(xc), C.double(yc),
		C.double(radius),
		C.double(angle1), C.double(angle2))
}

func (ctx *CairoContext) RelMoveTo(dx, dy float64) {
	C.cairo_rel_move_to(ctx.ptr, C.double(dx), C.double(dy))
}

func (ctx *CairoContext) RelLineTo(dx, dy float64) {
	C.cairo_rel_line_to(ctx.ptr, C.double(dx), C.double(dy))
}

func (ctx *CairoContext) RelCurveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64) {
	C.cairo_rel_curve_to(ctx.ptr,
		C.double(dx1), C.double(dy1),
		C.double(dx2), C.double(dy2),
		C.double(dx3), C.double(dy3))
}

func (ctx *CairoContext) Rectangle(x, y, width, height float64) {
	C.cairo_rectangle(ctx.ptr,
		C.double(x), C.double(y),
		C.double(width), C.double(height))
}

func (ctx *CairoContext) ClosePath() {
	C.cairo_close_path(ctx.ptr)
}

func (ctx *CairoContext) PathExtents() (left, top, right, bottom float64) {
	C.cairo_path_extents(ctx.ptr,
		(*C.double)(&left), (*C.double)(&top),
		(*C.double)(&right), (*C.double)(&bottom))
	return left, top, right, bottom
}

func (ctx *CairoContext) Paint() {
	C.cairo_paint(ctx.ptr)
}

func (ctx *CairoContext) PaintWithAlpha(alpha float64) {
	C.cairo_paint_with_alpha(ctx.ptr, C.double(alpha))
}

func (ctx *CairoContext) Mask(pattern Pattern) {
	C.cairo_mask(ctx.ptr, pattern.pattern)
}

func (ctx *CairoContext) Stroke() {
	C.cairo_stroke(ctx.ptr)
}

func (ctx *CairoContext) StrokePreserve() {
	C.cairo_stroke_preserve(ctx.ptr)
}

func (ctx *CairoContext) Fill() {
	C.cairo_fill(ctx.ptr)
}

func (ctx *CairoContext) FillPreserve() {
	C.cairo_fill_preserve(ctx.ptr)
}

func (ctx *CairoContext) CopyPage() {
	C.cairo_copy_page(ctx.ptr)
}

func (ctx *CairoContext) ShowPage() {
	C.cairo_show_page(ctx.ptr)
}

func (ctx *CairoContext) SetOperator(operator Operator) {
	C.cairo_set_operator(ctx.ptr, C.cairo_operator_t(operator))
}

func (ctx *CairoContext) SetSource(pattern *Pattern) {
	C.cairo_set_source(ctx.ptr, pattern.pattern)
}

func (ctx *CairoContext) SetSourceRGBA(color color.RGBA) {
	red := float64(color.R) / 255.0
	green := float64(color.G) / 255.0
	blue := float64(color.B) / 255.0
	aplha := float64(color.A) / 255.0
	C.cairo_set_source_rgba(ctx.ptr, C.double(red), C.double(green), C.double(blue), C.double(aplha))
}

func (ctx *CairoContext) SetSourceRGB(color color.RGBA) {
	red := float64(color.R) / 255.0
	green := float64(color.G) / 255.0
	blue := float64(color.B) / 255.0
	C.cairo_set_source_rgb(ctx.ptr, C.double(red), C.double(green), C.double(blue))
}

func (ctx *CairoContext) SetTolerance(tolerance float64) {
	C.cairo_set_tolerance(ctx.ptr, C.double(tolerance))
}

func (ctx *CairoContext) SetAntialias(antialias Antialias) {
	C.cairo_set_antialias(ctx.ptr, C.cairo_antialias_t(antialias))
}

func (ctx *CairoContext) SetFillRule(fill_rule FillRule) {
	C.cairo_set_fill_rule(ctx.ptr, C.cairo_fill_rule_t(fill_rule))
}

func (ctx *CairoContext) GetLineWidth() float64 {
	return float64(C.cairo_get_line_width(ctx.ptr))
}

func (ctx *CairoContext) SetLineWidth(width float64) {
	C.cairo_set_line_width(ctx.ptr, C.double(width))
}

func (ctx *CairoContext) SetLineCap(line_cap LineCap) {
	C.cairo_set_line_cap(ctx.ptr, C.cairo_line_cap_t(line_cap))
}

func (ctx *CairoContext) SetLineJoin(line_join LineJoin) {
	C.cairo_set_line_join(ctx.ptr, C.cairo_line_join_t(line_join))
}

func (ctx *CairoContext) SetDash(dashes []float64, num_dashes int, offset float64) {
	dashesp := (*C.double)(&dashes[0])
	C.cairo_set_dash(ctx.ptr, dashesp, C.int(num_dashes), C.double(offset))
}

func (ctx *CairoContext) SetMiterLimit(limit float64) {
	C.cairo_set_miter_limit(ctx.ptr, C.double(limit))
}

func (ctx *CairoContext) Translate(tx, ty float64) {
	C.cairo_translate(ctx.ptr, C.double(tx), C.double(ty))
}

func (ctx *CairoContext) Scale(sx, sy float64) {
	C.cairo_scale(ctx.ptr, C.double(sx), C.double(sy))
}

func (ctx *CairoContext) Rotate(angle float64) {
	C.cairo_rotate(ctx.ptr, C.double(angle))
}

func (ctx *CairoContext) SelectFontFace(name string, font_slant_t, font_weight_t int) {
	s := C.CString(name)
	C.cairo_select_font_face(ctx.ptr, s, C.cairo_font_slant_t(font_slant_t), C.cairo_font_weight_t(font_weight_t))
	C.free(unsafe.Pointer(s))
}

func (ctx *CairoContext) SetFontSize(size float64) {
	C.cairo_set_font_size(ctx.ptr, C.double(size))
}

func (ctx *CairoContext) ShowText(text string) {
	cs := C.CString(text)
	C.cairo_show_text(ctx.ptr, cs)
	C.free(unsafe.Pointer(cs))
}

func (ctx *CairoContext) TextPath(text string) {
	cs := C.CString(text)
	C.cairo_text_path(ctx.ptr, cs)
	C.free(unsafe.Pointer(cs))
}

func (ctx *CairoContext) TextExtents(text string) *TextExtents {
	cte := C.cairo_text_extents_t{}
	cs := C.CString(text)
	C.cairo_text_extents(ctx.ptr, cs, &cte)
	C.free(unsafe.Pointer(cs))
	te := &TextExtents{
		Xbearing: float64(cte.x_bearing),
		Ybearing: float64(cte.y_bearing),
		Width:    float64(cte.width),
		Height:   float64(cte.height),
		Xadvance: float64(cte.x_advance),
		Yadvance: float64(cte.y_advance),
	}
	return te
}

const (
	FONT_WEIGHT_NORMAL = iota
	FONT_WEIGHT_BOLD
)

type TextExtents struct {
	Xbearing float64
	Ybearing float64
	Width    float64
	Height   float64
	Xadvance float64
	Yadvance float64
}

type Operator int

const (
	OPERATOR_CLEAR = iota

	OPERATOR_SOURCE
	OPERATOR_OVER
	OPERATOR_IN
	OPERATOR_OUT
	OPERATOR_ATOP

	OPERATOR_DEST
	OPERATOR_DEST_OVER
	OPERATOR_DEST_IN
	OPERATOR_DEST_OUT
	OPERATOR_DEST_ATOP

	OPERATOR_XOR
	OPERATOR_ADD
	OPERATOR_SATURATE

	OPERATOR_MULTIPLY
	OPERATOR_SCREEN
	OPERATOR_OVERLAY
	OPERATOR_DARKEN
	OPERATOR_LIGHTEN
	OPERATOR_COLOR_DODGE
	OPERATOR_COLOR_BURN
	OPERATOR_HARD_LIGHT
	OPERATOR_SOFT_LIGHT
	OPERATOR_DIFFERENCE
	OPERATOR_EXCLUSION
	OPERATOR_HSL_HUE
	OPERATOR_HSL_SATURATION
	OPERATOR_HSL_COLOR
	OPERATOR_HSL_LUMINOSITY
)

const (
	FONT_SLANT_NORMAL = iota
	FONT_SLANT_ITALIC
	FONT_SLANT_OBLIQUE
)

type Antialias int

const (
	ANTIALIAS_DEFAULT Antialias = iota
	ANTIALIAS_NONE
	ANTIALIAS_GRAY
	ANTIALIAS_SUBPIXEL
)

type FillRule int

const (
	FILL_RULE_WINDING FillRule = iota
	FILL_RULE_EVEN_ODD
)

type LineCap int

const (
	LINE_CAP_BUTT LineCap = iota
	LINE_CAP_ROUND
	LINE_CAP_SQUARE
)

type LineJoin int

const (
	LINE_JOIN_MITER LineJoin = iota
	LINE_JOIN_ROUND
	LINE_JOIN_BEVEL
)

type Pattern struct {
	pattern *C.cairo_pattern_t
}

type PatternType int

const (
	PATTERN_TYPE_SOLID PatternType = iota
	PATTERN_TYPE_SURFACE
	PATTERN_TYPE_LINEAR
	PATTERN_TYPE_RADIAL
)
