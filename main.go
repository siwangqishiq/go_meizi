package main

import "fmt"

func main() {
	fmt.Println("Start go_meizi...")

	var crawler = NewCrawler()
	crawler.Load()
	crawler.FetchAblums()
}
