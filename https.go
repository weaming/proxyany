package main

import (
	"crypto/tls"
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
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}

func isHostAllowed(host string) bool {
	if strings.HasSuffix(host, allowedDomain) {
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
