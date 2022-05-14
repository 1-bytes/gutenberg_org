package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"spider/bootstrap"
	"spider/cmd"
	"spider/parser"
	"spider/pkg/fetcher"
	"strconv"
	"sync"
)

func main() {
	bootstrap.Setup()
	const BaseU = "https://www.gutenberg.org/ebooks/"

	runTimes := 68100
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(50, func(i interface{}) {
		u := BaseU + strconv.Itoa(i.(int))
		err := SpiderPage(u)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	})
	defer p.Release()

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(i)
	}
}

func SpiderPage(u string) error {
	fmt.Printf("request url: %s\n", u)
	resp, err := fetcher.Fetch(u)
	if err != nil {
		return fmt.Errorf("error: send request failed: %s", err)
	}
	p := parser.Parser{}
	p.GetDetail(u, resp)
	cmd.SavaData(p)
	return nil
}
