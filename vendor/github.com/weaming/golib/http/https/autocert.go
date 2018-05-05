package https

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type Config struct {
	CacheDir         string
	HTTPSecureServer *http.Server
	HTTPServer       *http.Server
	IsHostAllowed    func(host string) bool
	Manager          *autocert.Manager
}

func (c *Config) ListenAndServeTLS() error {
	hostPolicy := func(ctx context.Context, host string) error {
		if c.IsHostAllowed(host) {
			return nil
		}
		return fmt.Errorf("host %v is not allowed", host)
	}

	c.Manager = &autocert.Manager{
		Cache:      autocert.DirCache(c.CacheDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
	}
	c.HTTPSecureServer.TLSConfig = &tls.Config{GetCertificate: c.Manager.GetCertificate}

	// https
	go func() {
		c.HTTPSecureServer.Addr = ":https"
		err := c.HTTPSecureServer.ListenAndServeTLS("", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// http
	c.HTTPServer.Handler = c.Manager.HTTPHandler(c.HTTPServer.Handler)
	c.HTTPServer.Addr = ":http"
	return c.HTTPServer.ListenAndServe()
}
