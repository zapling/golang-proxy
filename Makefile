.PHONY: build start

build:
	go build -o build/golang-proxy

start:
	build/./golang-proxy 5000 http://localhost:80