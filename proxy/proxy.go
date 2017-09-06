package main

import (
	"net/http"
	"net/http/httputil"
)

var tableName = "SerialHostLookup"
var region = "us-east-1"

func handler(w http.ResponseWriter, r *http.Request) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8081"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/receiveMsg", handler)
	panic(http.ListenAndServe(":8080", nil))
}
