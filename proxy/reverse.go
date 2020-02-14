package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy http server
type ReverseProxy struct {
	Port   string
	Target string
}

// Serve start the http server
func (p *ReverseProxy) Serve() {
	p.setup()
	fmt.Printf(
		"\nServing reverse proxy server...\nhttp://localhost:%v -> %v\n",
		p.Port, p.Target,
	)

	if err := http.ListenAndServe(":"+p.Port, nil); err != nil {
		panic(err)
	}
}

func (p *ReverseProxy) setup() {
	http.HandleFunc("/", p.handleRequest)
}

func (p *ReverseProxy) logRequest(req *http.Request) {
	fmt.Printf(
		"[%v] {%v} %v %v %v\n",
		req.Method,
		req.RemoteAddr,
		req.URL,
		req.Header["Content-Type"],
		req.ContentLength,
	)
}

func (p *ReverseProxy) handleRequest(res http.ResponseWriter, req *http.Request) {
	p.logRequest(req)
	url, _ := url.Parse(p.Target)
	p.setHTTPSHeaders(url, req)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func (p *ReverseProxy) setHTTPSHeaders(url *url.URL, req *http.Request) {
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("host"))
	req.Host = url.Host
}
