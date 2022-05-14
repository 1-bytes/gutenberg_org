package main

import (
	"fmt"
	"log"
	"spider/bootstrap"
	"spider/cmd"
	"spider/parser"
	"spider/pkg/fetcher"
	"strconv"
)

func main() {
	bootstrap.Setup()
	baseU := "https://www.gutenberg.org/ebooks/"
	//for i := 1; i < 68100; i++ {
	for i := 8; i < 9; i++ {
		u := baseU + strconv.Itoa(i)
		fmt.Printf("request url: %s\n", u)
		resp, err := fetcher.Fetch(u)
		if err != nil {
			log.Printf("error: send request failed: %s", err)
			continue
		}
		p := parser.Parser{}
		p.GetDetail(u, resp)
		cmd.SavaData(&p)

		//fmt.Println(p.Title)
		//fmt.Println(p.Author)
		//fmt.Println(p.Language)
		//fmt.Println(p.ReleaseDate)
		//fmt.Println(p.DownloadCount)
		//fmt.Println(p.CoverImage)
		//fmt.Println(p.Content)
		//fmt.Println("\n")
	}
}
