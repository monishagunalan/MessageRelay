package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	fmt.Fprintf(w, "\nRequest msg:  %s", string(body))
}

func main() {
	http.HandleFunc("/receiveMsg", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
