package cmd

import (
	"bbc_com/bootstrap"
	"bbc_com/parser"
	"bbc_com/pkg/config"
	elasticsearch "bbc_com/pkg/elastic"
	"bbc_com/pkg/queued"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/olivere/elastic/v7"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// NewCollector 传入配置信息，创建并返回一个 colly 的 collector 实例
func NewCollector(options ...colly.CollectorOption) *colly.Collector {
	c := colly.NewCollector(options...)
	// 代理设置
	proxyAddress := config.GetString("spider.socks5")
	rp, err := proxy.RoundRobinProxySwitcher(proxyAddress)
	if err != nil {
		log.Println("attempt to use Socks5 proxy failed.")
		panic(err)
	}
	if proxyAddress == "" {
		rp = nil
	}

	// 爬虫速度以及响应时间等参数的控制
	c.WithTransport(&http.Transport{
		Proxy: rp,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	})
	// 初始化 Redis Storage，将其用作爬虫的持久化队列
	if err := c.SetStorage(bootstrap.Storage); err != nil {
		panic(err)
	}
	return c
}

// SpiderCallbacks colly 的回调函数
func SpiderCallbacks(c *colly.Collector) {
	// 请求发起之前要处理的一些事件
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
		r.Headers.Set("Referer", "https://www.bbc.com")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	})

	// 抓取新的页面
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		_ = queued.Queued.AddURL(url)
	})

	// 处理请求结果
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		domain := r.Request.URL.Host

		articleData, err := parser.GetArticleAllData(r.Body)
		if err != nil {
			return
		}
		title := parser.Title(&articleData)
		author := parser.Author(&articleData)
		category := parser.Category(url)
		paragraphs, err := parser.Content(&articleData)
		if err != nil {
			log.Printf("Error: %s\n", err)
		}
		if len(paragraphs) == 0 {
			return
		}
		// 插入到MySQL
		model := parser.DictArticleModel{
			Type:                parser.TypeMap[category],
			Title:               title,
			Author:              author,
			MostRecentlyUpdated: "",
			SourceDomain:        4,
		}

		err = SaveDataToMySQL("dict_article", &model)
		if err != nil {
			log.Printf("MySQL SaveData error: %v\n", err)
			return
		}

		for _, paragraph := range paragraphs {
			//log.Printf("ID: %d\n", id)
			//log.Printf("Title: %s\n", title)
			//log.Printf("Author: %s\n", author)
			//log.Printf("Category: %s\n", category)
			//log.Printf("ReleaseDate: %s\n", releaseDate)
			//log.Printf("EN: %s\n", paragraph["EN"])
			//log.Printf("CN: %s\n", paragraph["CN"])
			//log.Println()

			data := parser.JsonData{
				ID:           strconv.Itoa(model.ID),
				SourceDomain: domain,
				SourceURL:    url,
				Paragraph:    paragraph,
			}
			if err = SaveDataToElastic("dict_article", "", &data); err != nil {
				log.Printf("ElasticSearch SaveData error: %v\n", err)
			}
		}
	})

	// 错误处理
	c.OnError(func(resp *colly.Response, err error) {
		//err = resp.Request.Retry()
		err = queued.Queued.AddRequest(resp.Request)
		if err != nil {
			log.Println("Request URL:", resp.Request.URL, "failed with response:", resp, "\nError:", err)
		}
	})
}

// SaveDataToElastic 存储数据至 ES
func SaveDataToElastic(index string, id string, data *parser.JsonData) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var e *elastic.IndexService
	e = elasticsearch.GetInstance().Index()
	if id != "" {
		e.Id(id)
	}
	_, err = e.Index(index).
		BodyJson(string(j)).
		Timeout("5s").
		Do(context.Background())
	if err != nil {
		fmt.Println(id, data)
		return err
		//log.Printf("%+v: %+v\n", do.Result, do.Id)
	}
	return err
}

// SaveDataToMySQL 存储数据至 mysql
func SaveDataToMySQL(tables string, data *parser.DictArticleModel) error {
	db := bootstrap.DB
	tx := db.Table(tables).Create(data)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}
