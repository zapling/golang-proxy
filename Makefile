.PHONY: build start install

build:
	go build -o build/golang-proxy

start:
	build/./golang-proxy 5000 http://localhost:80

install:
	go install