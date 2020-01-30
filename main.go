package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func printUsage() {
	text := "golang-proxy [PORT] [TARGET]\n" +
		"Example: golang-proxy 5000 http://localhost:80\n\n" +
		"PORT:   The port that the proxy should run on.\n" +
		"TARGET: The target address that the proxy should redirect requests to."

	fmt.Println(text)
}

func exitError(message string) {
	printUsage()
	fmt.Printf("\n---\nError: %v\n", message)
	os.Exit(1)
}

func getParams(args []string) (string, string) {
	var port string
	var target string

	for i := 1; i < len(args); i++ {
		switch i {
		case 1:
			port = args[i]
		case 2:
			target = args[i]
		}
	}

	if port == "" {
		exitError("PORT most be specified")
	}

	if target == "" {
		exitError("TARGET mot be specified")
	}

	return port, target
}

func printRequestInfo(req *http.Request) {
	fmt.Printf(
		"[%v] {%v} %v %v %v\n",
		req.Method,
		req.RemoteAddr,
		req.URL,
		req.Header["Content-Type"],
		req.ContentLength,
	)
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	printRequestInfo(req)

	_, target := getParams(os.Args)
	url, _ := url.Parse(target)

	// https redirects
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("host"))
	req.Host = url.Host

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func main() {
	port, target := getParams(os.Args)

	fmt.Println("Server http server...")
	fmt.Printf("http://localhost:%v --> %v\n", port, target)

	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
