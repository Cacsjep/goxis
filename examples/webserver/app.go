package main

import (
	"embed"
	"net/http"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed index.html
var f embed.FS

//go:embed static/*
var embedDirStatic embed.FS

// This example demonstrates how to create a reverse proxy webserver with fiber.
//
// To use a webserver like fiber we use the reverse proxy support for ACAP,
// therefrore we need to declare reverseProxy item in Manifest.
//
// Currently there is now way to declare a settingsPage directly to the reverseProxy path,
// so we need a redirect.html that redirects to the correct path.
// Thats needed because we want serve our own html files.
//
// Original C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/web-server
func main() {
	var err error
	var app *acapapp.AcapApplication
	var baseUri string

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app = acapapp.NewAcapApplication()

	// Fiber
	fapp := fiber.New()
	if baseUri, err = app.AcapWebBaseUri(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Index.html hosting
	fapp.Use(baseUri, filesystem.New(filesystem.Config{
		Root: http.FS(f),
	}))

	// Static file hosting
	fapp.Use(baseUri+"/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
	}))

	// Start the webserver
	fapp.Listen("127.0.0.1:2001")

	app.Close()
}
