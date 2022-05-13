package main

import (
	"bbc_com/parser"
	"bbc_com/pkg/fetcher"
	"fmt"
	"log"
	"strconv"
)

func main() {
	baseU := "https://www.gutenberg.org/ebooks/"
	//for i := 1; i < 68100; i++ {
	for i := 67998; i < 67999; i++ {
		u := baseU + strconv.Itoa(i)
		fmt.Printf("request url: %s\n", u)
		bytes, err := fetcher.Fetch(u)
		if err != nil {
			log.Printf("error: send request faild: %s", err)
			continue
		}
		p := parser.Parser{}
		p.GetDetail(bytes)

		fmt.Println(p.Title)
		fmt.Println(p.Author)
		fmt.Println(p.Language)
		fmt.Println(p.ReleaseDate)
		fmt.Println(p.DownloadCount)
		fmt.Println(p.CoverImage)
		fmt.Println(p.Content)
		fmt.Println("\n")
	}
}
