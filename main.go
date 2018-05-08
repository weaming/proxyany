package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/weaming/proxyany/reverseproxy"
)

var (
	bind    = ":20443"
	https   = false
	cfgPath = "config.json"
	version = "version 1.2"
	mg      *reverseproxy.MapGroup
)

func init() {
	fmt.Println(version)
	flag.StringVar(&cfgPath, "config", cfgPath, "file path domain mapping config in json format")
	flag.StringVar(&bind, "bind", bind, "local bind [<host>]:<port>")
	flag.BoolVar(&https, "https", https, "HTTPS mode, auto certification from let's encrypt")
	flag.Parse()

	mg = reverseproxy.LoadMapGroupFromJson(cfgPath)
}

func main() {
	srv := newProxyServer()
	if https {
		config := NewConfig(srv)
		config.IsHostAllowed = isHostAllowed

		fmt.Printf("listening :443\n")
		config.ListenAndServeTLS()
	} else {
		srv.Addr = bind

		fmt.Printf("listening %v\n", bind)
		err := srv.ListenAndServe()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func newProxyServer() *http.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := reverseproxy.NewReverseProxy(mg)
		proxy.ServeHTTP(w, r)
	})

	srv := NewHTTPServer()
	srv.Handler = handler
	return srv
}

func isHostAllowed(host string) bool {
	if mg.GetMapping(host) != nil {
		return true
	}
	return false
}
