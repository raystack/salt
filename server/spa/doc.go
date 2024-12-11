/*
Package spa provides a simple and efficient HTTP handler for serving
Single Page Applications (SPAs).

The handler serves static files from an embedded file system and falls
back to serving an index file for client-side routing. Optionally, it
supports gzip compression for optimizing responses.

Features:
  - Serves static assets from an embedded file system.
  - Fallback to an index file for client-side routing.
  - Optional gzip compression for supported clients.

Usage:

To use this package, embed your SPA's build assets into your binary using
the `embed` package. Then, create an SPA handler using the `Handler` function
and register it with an HTTP server.

Example:

	package main

	import (
		"embed"
		"log"
		"net/http"

		"yourmodule/spa"
	)

	//go:embed build/*
	var build embed.FS

	func main() {
		handler, err := spa.Handler(build, "build", "index.html", true)
		if err != nil {
			log.Fatalf("Failed to initialize SPA handler: %v", err)
		}

		log.Println("Serving SPA on http://localhost:8080")
		http.ListenAndServe(":8080", handler)
	}
*/
package spa
