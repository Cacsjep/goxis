package main

import (
	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

var (
	err          error
	subscription int
	app          *goxis.AcapApplication
)

// This example uses axoverlay example to draw a rectangle
func main() {
	if !acap.AxOverlayIsBackendSupported(acap.AxOverlayCairoImageBackend) {
		panic("cairo not supported")
	}

	//  Initialize the library
	//  Setup colors
	// Get max resolution for width and height
	// Create a large overlay using Palette color space
	// Create an text overlay using ARGB32 color space
	// Draw overlays
	// Start animation timer

	// Enter loop

	// Destroy overlay
	// cleanup
}
