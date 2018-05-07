package main

import (
	"crypto/tls"
	"net/http"
	"time"

	libhttps "github.com/weaming/golib/http/https"
)

func NewHTTPServer() *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), //disable http2
	}
}

func NewConfig(srv *http.Server) *libhttps.Config {
	config := &libhttps.Config{
		CacheDir:         ".",
		HTTPSecureServer: srv,
		HTTPServer:       libhttps.NewRedirectServer(),
		IsHostAllowed:    nil,
	}
	return config
}
