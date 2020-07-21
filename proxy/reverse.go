package proxy

import (
	"fmt"
    "time"
    "strings"
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
    localhost_addr := "[::1]"
    remote_addr := req.RemoteAddr

    if strings.Contains(remote_addr, localhost_addr) {
        remote_addr = "127.0.0.1:" + remote_addr[6:]
    }

	fmt.Printf(
        "%v - [%v] %v [L:%v]\n",
        time.Now().Format(time.RFC3339),
        req.Method,
		remote_addr,
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
