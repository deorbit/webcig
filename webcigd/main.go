// webcigd starts the webcig HTTP server listening for requests.
package main

import (
	"github.com/deorbit/webcig/server"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("WEBCIGPORT")
	webTemplateDir := os.Getenv("WEBCIG_TEMPLATE_DIR")
	webStaticDir := os.Getenv("WEBCIG_STATIC_DIR")

	if port == "" {
		port = "8081"
	}

	http.ListenAndServe(":"+port, server.New(webTemplateDir, webStaticDir))
}
