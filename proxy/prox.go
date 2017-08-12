package prox

import (
	"net/http"

	"net/http/httputil"
	"net/url"
)

/*Prox the structure for the proxy server*/
type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

/*New creates a new Prox structure*/
func New(target string) *Prox {
	url, _ := url.Parse(target)

	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

/*Handle generic handler for proxied calls*/
func (p *Prox) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Go-Proxy", "Prox")

	p.proxy.ServeHTTP(w, r)
}
