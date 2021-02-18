package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", wsEndpoint)

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go cleanCache()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	log.Println("Starting web server on port 7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}
