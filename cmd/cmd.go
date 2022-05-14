package cmd

import (
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"spider/bootstrap"
	"spider/model"
	"spider/parser"
	"spider/pkg/elastic"
	"strconv"
	"strings"
)

func SavaData(bookDetail *parser.Parser) {
	db := bootstrap.DB
	splitParagraphs := SplitParagraph(&bookDetail.Content)
	// 存储至 MySQL
	articleModel := model.DictArticle{
		Type:                "1",
		Title:               bookDetail.Title,
		Author:              bookDetail.Author,
		CoverImage:          bookDetail.CoverImage,
		Category:            "0",
		DownloadCount:       bookDetail.DownloadCount,
		ReleaseDate:         bookDetail.ReleaseDate,
		MostRecentlyUpdated: "",
		SourceDomain:        "5",
		Status:              "1",
		Rank:                99,
	}
	db.Table(articleModel.TableName()).Create(&articleModel)
	// 存储至 ElasticSearch
	for _, paragraph := range splitParagraphs {
		articleID := int(articleModel.ID)
		data := parser.ElasticSearchData{
			ID:           strconv.Itoa(articleID),
			SourceDomain: "www.gutenberg.org",
			SourceUrl:    bookDetail.URL,
			Paragraph: struct {
				EN string `json:"en"`
			}{
				EN: paragraph,
			},
		}
		dataJson, err := json.Marshal(data)
		if err != nil {
			log.Printf("error: elasticsearch json marshal failed: %s", err)
			continue
		}

		client := elastic.GetInstance().Index()
		do, err := client.
			Index("dict_article").
			BodyJson(string(dataJson)).
			Timeout("10s").
			Do(context.Background())
		if err != nil {
			log.Printf("error: save data to the elasticsearch failed: %s", err)
			continue
		}
		// 存储至 MySQL
		articleParagraph := model.DictArticleParagraph{
			ElasticID:     do.Id,
			ArticleID:     articleID,
			ByteCount:     len(paragraph),
			Fre:           "",
			Fkgl:          "",
			SchoolLvClass: 0,
			SchoolLvlName: "",
			TechWordLv:    "",
			CefrWordLv:    "",
			Status:        1,
		}
		db.Table(articleParagraph.TableName()).Create(&articleParagraph)
	}
}

// SplitParagraph 段落拆分
func SplitParagraph(resp *string) []string {
	paragraphs := strings.Split(*resp, "\r\n\r\n")
	var result []string
	for _, paragraph := range paragraphs {
		paragraph = strings.ReplaceAll(paragraph, "\r\n", " ")
		paragraph = strings.TrimSpace(paragraph)
		paragraph = strings.Trim(paragraph, "***")
		if paragraph == "" {
			continue
		}
		result = append(result, paragraph)
	}
	return result
}
