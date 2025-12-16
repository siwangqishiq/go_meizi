package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start go_meizi...")

	var crawler = NewCrawler()
	crawler.Load()
	crawler.FetchAblums()

	go func ()  {
		for {
			time.Sleep(4 * time.Hour)
			fmt.Println("Start crawler Task...")
			crawler.FetchAblums()
		}
	}()

	var httpServer = NewHttpServer(crawler)
	httpServer.StartServer()
}
