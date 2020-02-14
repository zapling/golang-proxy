package main

import (
	"fmt"
	"os"

	"github.com/zapling/golang-proxy/proxy"
)

var version = "1.0.2"

func usage() {
	fmt.Println(
		"golang-proxy [PORT] [TARGET]\nVersion " + version + "\nhttps://github.com/zapling/golang-proxy/\n\n" +
			"Example: golang-proxy 5000 http://localhost:80\n\n" +
			"PORT:   The port that the proxy should run on.\n" +
			"TARGET: The target address that the proxy should redirect requests to.",
	)
}

func getArgs() (string, string) {
	var port, target string

	args := os.Args

	if len(args) < 3 {
		usage()
		os.Exit(0)
	}

	port = args[1]
	target = args[2]

	if port == "" || target == "" {
		fmt.Println("Insufficient arguments")
		os.Exit(1)
	}

	return port, target
}

func main() {
	port, target := getArgs()
	p := proxy.ReverseProxy{Port: port, Target: target}
	p.Serve()
}
