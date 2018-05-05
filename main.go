package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/cssivision/reverseproxy"
)

var (
	target        = "https://www.google.com"
	bind          = ":20443"
	https         = false
	allowedDomain = "bitsflow.org"
)

func init() {
	flag.StringVar(&target, "from", target, "your reverse proxy target url, including path is allowed, then your visit path will be append to it")
	flag.StringVar(&bind, "to", bind, "local bind [<host>]:<port>")
	flag.StringVar(&allowedDomain, "domain", allowedDomain, "domain allowed to access, all sub domains will be allowed too")
	flag.BoolVar(&https, "https", https, "HTTPS mode, auto certification from let's encrypt")

	flag.Parse()
}

func main() {
	srv := newProxyServer()
	if https {
		config := newConfig(srv)
		fmt.Printf("proxy from %v to :443\n", target)
		config.ListenAndServeTLS()
	} else {
		srv.Addr = bind

		fmt.Printf("proxy from %v to %v\n", target, bind)
		err := srv.ListenAndServe()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func newProxyServer() *http.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		proxy := reverseproxy.NewReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})

	srv := newHTTPServer()
	srv.Handler = handler
	return srv
}
