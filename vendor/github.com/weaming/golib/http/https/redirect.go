package https

import (
	"net/http"
)

func RedirectHTTPS(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusFound)
}

func NewRedirectServer() *http.Server {
	return &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(RedirectHTTPS),
	}
}
