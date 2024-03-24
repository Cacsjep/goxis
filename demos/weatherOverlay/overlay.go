package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/Cacsjep/goxis/pkg/acap"
)

type Compass struct {
	WindDirectionDegrees float64
	Textsize             float64
	CenterX              float64
	CenterY              float64
	LabelPosition        float64
	Radius               float64
	Theme                CompassTheme
	DirectionLineWidth   float64
	BaseLineWidth        float64
	DegreeStep           float64
	NordDirection        float64
	DegreeTextOffset     float64
	DegreeLineOffset     float64
}

type CompassTheme struct {
	DegreeDirectionLine color.RGBA
	DegreeLines         color.RGBA
	DegreeStep          color.RGBA
	Circle              color.RGBA
	Text                color.RGBA
}

func adjustmentCallback(adjustmentEvent *acap.OverlayAdjustmentEvent) {
	*adjustmentEvent.OverlayWidth = adjustmentEvent.Stream.Width
	*adjustmentEvent.OverlayHeight = adjustmentEvent.Stream.Height
}

func renderCallback(renderEvent *acap.OverlayRenderEvent) {
	wapp := renderEvent.Userdata.(*WeatherApp)
	if renderEvent.OverlayId == wapp.TemperatureOverlayId {

		c := Compass{
			NordDirection:        wapp.NordDirection,
			Radius:               wapp.Size,
			LabelPosition:        1.2,
			Textsize:             18,
			WindDirectionDegrees: float64(wapp.LastData.CurrentWeather.Winddirection),
			DegreeStep:           10,
			DegreeTextOffset:     0.86,
			DegreeLineOffset:     0.97,
			Theme: CompassTheme{
				Circle:              wapp.CircleColor,
				Text:                wapp.CircleColor,
				DegreeDirectionLine: wapp.Color,
				DegreeLines:         wapp.Color,
				DegreeStep:          wapp.Color,
			},
			DirectionLineWidth: 2,
			BaseLineWidth:      4,
		}

		switch wapp.Position {
		case acap.AxOverlayTopLeft:
			c.CenterX = c.Radius*c.LabelPosition + c.Textsize
			c.CenterY = c.Radius*c.LabelPosition + c.Textsize
		case acap.AxOverlayTopRight:
			c.CenterX = float64(renderEvent.Stream.Width) - c.Radius*c.LabelPosition + c.Textsize
			c.CenterY = c.Radius*c.LabelPosition + c.Textsize
		case acap.AxOverlayBottomLeft:
			c.CenterX = c.Radius*c.LabelPosition + c.Textsize
			c.CenterY = float64(renderEvent.Stream.Height) - c.Radius*c.LabelPosition - c.Textsize
		case acap.AxOverlayBottomRight:
			c.CenterX = float64(renderEvent.Stream.Width) - c.Radius*c.LabelPosition - c.Textsize
			c.CenterY = float64(renderEvent.Stream.Height) - c.Radius*c.LabelPosition - c.Textsize
		}

		// Draw the base and points of the compass
		DrawCompassBase(renderEvent.CairoCtx, c)
		DrawDegreeLines(renderEvent.CairoCtx, c)
		DrawCompassPoints(renderEvent.CairoCtx, c)
		HighlightDirection(renderEvent.CairoCtx, c)
		DrawText(renderEvent.CairoCtx, c, wapp)
	}
}

func DrawCompassBase(ctx *acap.CairoContext, compass Compass) {
	ctx.NewPath()
	ctx.Arc(compass.CenterX, compass.CenterY, compass.Radius, 0, 2*math.Pi)
	ctx.SetSourceRGB(compass.Theme.Circle)
	ctx.SetLineWidth(compass.BaseLineWidth)
	ctx.Stroke()
}

func DrawDegreeLines(ctx *acap.CairoContext, compass Compass) {
	for degree := 0.0; degree < 360.0; degree += compass.DegreeStep {
		adjustedDegree := degree + compass.NordDirection
		adjustedDegree -= 90
		if adjustedDegree >= 360 {
			adjustedDegree -= 360
		} else if adjustedDegree < 0 {
			adjustedDegree += 360
		}

		// Calculate the outer point of the line on the circle
		radOuter := adjustedDegree * (math.Pi / 180)
		outerX := compass.CenterX + compass.Radius*math.Cos(radOuter)
		outerY := compass.CenterY + compass.Radius*math.Sin(radOuter)

		// Calculate the inner point of the line, making the line go towards the center
		innerX := compass.CenterX + (compass.Radius*compass.DegreeLineOffset)*math.Cos(radOuter)
		innerY := compass.CenterY + (compass.Radius*compass.DegreeLineOffset)*math.Sin(radOuter)

		// Draw the line
		ctx.MoveTo(outerX, outerY)
		ctx.LineTo(innerX, innerY)

		// Set the color and width of the degree lines
		ctx.SetSourceRGBA(compass.Theme.DegreeStep)
		ctx.SetLineWidth(1) // Adjust the line width as needed
		ctx.Stroke()

		// Degree labels
		textRadius := compass.Radius * compass.DegreeTextOffset // Adjust this value as needed
		textX := compass.CenterX + textRadius*math.Cos(radOuter)
		textY := compass.CenterY + textRadius*math.Sin(radOuter)

		// Format the degree text. Adjust the string format as needed.
		degreeText := fmt.Sprintf("%d", int(degree))

		// Set font for degree text
		ctx.SetFontSize(compass.Textsize - (compass.Textsize / 2)) // Smaller font size for the degree text

		// Calculate text width and height to adjust positioning
		extents := ctx.TextExtents(degreeText)
		textWidth := extents.Width
		textHeight := extents.Height

		// Adjust text position to center it around the calculated point
		ctx.MoveTo(textX-textWidth/2, textY+textHeight/2)

		// Draw the text
		ctx.ShowText(degreeText)
	}
}

func isWholeNumber(degree float64) bool {
	return math.Mod(degree, 1) == 0
}

func DrawCompassPoints(ctx *acap.CairoContext, compass Compass) {
	cardinalPoints := map[string]float64{
		"N":   0, // Nord
		"NNO": 22.5,
		"NO":  45, // Nordost
		"ONO": 67.5,
		"O":   90, // Ost
		"OSO": 112.5,
		"SO":  135, // Südost
		"SSO": 157.5,
		"S":   180, // Süd
		"SSW": 202.5,
		"SW":  225, // Südwest
		"WSW": 247.5,
		"W":   270, // West
		"WNW": 292.5,
		"NW":  315, // Nordwest
		"NNW": 337.5,
	}

	ctx.SetFontSize(compass.Textsize)
	ctx.SetSourceRGB(compass.Theme.Text) // Assuming compass.Theme.Text is defined elsewhere
	for point, angle := range cardinalPoints {
		// Apply NordDirection adjustment
		adjustedAngle := angle - 90 + compass.NordDirection
		// Normalize the angle to 0-360 range
		if adjustedAngle >= 360 {
			adjustedAngle -= 360
		} else if adjustedAngle < 0 {
			adjustedAngle += 360
		}
		// Convert the adjusted angle to radians
		rad := adjustedAngle * (math.Pi / 180)
		labelPos := compass.LabelPosition

		if isWholeNumber(angle) {
			ctx.SelectFontFace("serif", 0, acap.FONT_WEIGHT_BOLD)
			ctx.SetFontSize(compass.Textsize)
		} else {
			ctx.SelectFontFace("serif", 0, 0)
			ctx.SetFontSize(compass.Textsize - (compass.Textsize / 3))
			labelPos = labelPos - (compass.LabelPosition / 2)
		}

		// Calculate the position for the text
		textX := compass.CenterX + (compass.Radius*labelPos)*math.Cos(rad)
		textY := compass.CenterY + (compass.Radius*labelPos)*math.Sin(rad)

		extents := ctx.TextExtents(point)
		// Adjust the positioning of the text based on its extents
		ctx.MoveTo(textX-extents.Width/2, textY+extents.Height/2)

		ctx.ShowText(point)
	}
}

func HighlightDirection(ctx *acap.CairoContext, compass Compass) {
	adjustedWindDirection := compass.WindDirectionDegrees + compass.NordDirection
	adjustedWindDirection -= 90
	if adjustedWindDirection >= 360 {
		adjustedWindDirection -= 360
	} else if adjustedWindDirection < 0 {
		adjustedWindDirection += 360
	}

	rad := adjustedWindDirection * (math.Pi / 180)
	endX := compass.CenterX + compass.Radius*math.Cos(rad)
	endY := compass.CenterY + compass.Radius*math.Sin(rad)
	ctx.SetSourceRGB(compass.Theme.DegreeDirectionLine)
	ctx.MoveTo(compass.CenterX, compass.CenterY)
	ctx.LineTo(endX, endY)
	ctx.SetLineWidth(compass.DirectionLineWidth)
	ctx.Stroke()
}

func DrawText(ctx *acap.CairoContext, compass Compass, wapp *WeatherApp) {
	text := fmt.Sprintf(
		"%.1f%s",
		wapp.LastData.CurrentWeather.Windspeed,
		wapp.LastData.CurrentWeatherUnits.Windspeed,
	)
	ctx.SelectFontFace("serif", 0, acap.FONT_WEIGHT_BOLD)
	ctx.SetFontSize(compass.Textsize)
	ctx.SetSourceRGB(compass.Theme.Text)

	// Calculate text width and height to center it
	extents := ctx.TextExtents(text)
	textWidth := extents.Width
	textHeight := extents.Height

	startX := compass.CenterX - (textWidth / 2)
	startY := compass.CenterY + (textHeight / 2)

	// Move to calculated start position and show text
	ctx.MoveTo(startX, startY-(compass.Textsize*2))
	ctx.ShowText(text)

	temp_text := fmt.Sprintf(
		"%.1f%s",
		wapp.LastData.CurrentWeather.Temperature,
		wapp.LastData.CurrentWeatherUnits.Temperature,
	)

	extents = ctx.TextExtents(temp_text)
	textWidth = extents.Width
	textHeight = extents.Height

	startX = compass.CenterX - (textWidth / 2)
	startY = compass.CenterY + (textHeight / 2)

	ctx.MoveTo(startX, startY+(compass.Textsize*2))
	ctx.ShowText(temp_text)
}

func (w *WeatherApp) Redraw() {
	if err := w.OvProvider.Redraw(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to redraw overlays: %s", err.Error())
	}
}
