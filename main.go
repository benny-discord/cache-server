package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", wsEndpoint)

	var port = 7000

	flag.IntVar(&port, "p", 7000, "the port to run the web server on")
	flag.Parse()

	var portstring = fmt.Sprintf(":%d", port)

	log.Printf("Starting web server on port %d\n", port)
	log.Fatal(http.ListenAndServe(portstring, nil))
}
