package main

import (
	"net/http"
	"strings"
	"time"

	libhttps "github.com/weaming/golib/http/https"
)

func newHTTPServer() *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func isHostAllowed(host string) bool {
	if strings.HasSuffix(host, "bitsflow.org") {
		return true
	}
	return false
}

func newConfig(srv *http.Server) *libhttps.Config {
	config := &libhttps.Config{
		CacheDir:         ".",
		HTTPSecureServer: srv,
		HTTPServer:       libhttps.NewRedirectServer(),
		IsHostAllowed:    isHostAllowed,
	}
	return config
}
