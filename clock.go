package main

import "time"

func init() {
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
}
