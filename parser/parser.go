package parser

import (
	"fmt"
	"log"
	"regexp"
	"spider/pkg/downloader"
	"spider/pkg/fetcher"
	"strconv"
	"strings"
)

var content string

type Parser struct {
	URL           string
	ID            string
	Author        string
	Title         string
	Language      string
	ReleaseDate   string
	DownloadCount int
	CoverImage    string
	Content       string
	response      []byte
}

var (
	authorRe        = regexp.MustCompile(`<a href="/ebooks/author[^>]*?">([^<]*?)</a></td>`)
	titleRe         = regexp.MustCompile(`<meta name="title" content="([^"]*?)">`)
	languageRe      = regexp.MustCompile(`<th>Language</th>\s<td>(<a href=[^>]*?>)?([^<]*?)(</a>)?</td>`)
	releaseDateRe   = regexp.MustCompile(`<th>Release Date</th>\s<td[^>]*?>([^<]*?)</td>`)
	downloadCountRe = regexp.MustCompile(`<td itemprop="interactionCount">(\d*?) downloads in`)
	coverImageRe    = regexp.MustCompile(`<img class="cover-art" src="([^"]*?)"\s`)
	contentURLRe    = regexp.MustCompile(`<a href="([^"]*?)" type="text/plain`)
	contentRe       = regexp.MustCompile(
		`\*\*\*\sSTART\sOF\s(THE|THIS)\sPROJECT.+\*\*\*([\s\S]+)\*\*\*\sEND\sOF\s(THE|THIS)\sPROJECT.+\*\*\*`)
)

func (p *Parser) GetDetail(url string, resp []byte) {
	p.response = resp
	var err error
	p.URL = url
	p.Author = p.regexpMatch(authorRe, 1, p.response)
	p.Title = p.regexpMatch(titleRe, 1, p.response)
	p.Language = p.regexpMatch(languageRe, 1, p.response)
	p.ReleaseDate = p.regexpMatch(releaseDateRe, 1, p.response)
	p.DownloadCount = p.downloadCount()
	if p.CoverImage, err = p.coverImage(); err != nil {
		log.Println(err)
	}
	if p.Content, err = p.content(); err != nil {
		log.Println(err)
	}
}

// downloadCount 获取下载次数
func (p *Parser) downloadCount() int {
	c, _ := strconv.Atoi(p.regexpMatch(downloadCountRe, 1, p.response))
	return c
}

// coverImage 封面图片
func (p *Parser) coverImage() (string, error) {
	idRe := regexp.MustCompile(`/epub/([^/]*?)/`)
	cover := p.regexpMatch(coverImageRe, 1, p.response)
	if cover == "" {
		return "", fmt.Errorf("cover image not found, url: %s\n", p.URL)
	}
	id := p.regexpMatch(idRe, 1, []byte(cover))
	err := downloader.DownloadImage(cover, "files/cover_image/"+id+".jpg")
	if err != nil {
		return "", fmt.Errorf("download cover image failed: %s\n", err)
	}
	return id + ".jpg", nil
}

// content 图书内容
func (p *Parser) content() (string, error) {
	contentURL := p.regexpMatch(contentURLRe, 1, p.response)
	if contentURL == "" {
		return "", fmt.Errorf("books url not found")
	}
	contentURL = "https://www.gutenberg.org" + contentURL
	bytes, err := fetcher.Fetch(contentURL)
	if err != nil {
		return "", fmt.Errorf("fetch books content failed: %s", err)
	}
	c := p.regexpMatch(contentRe, 2, bytes)
	if c == "" {
		return "", fmt.Errorf("books content not found")
	}
	return c, nil
}

// regexpMatch 通用正则匹配方法
func (p *Parser) regexpMatch(re *regexp.Regexp, index int, content []byte) string {
	match := re.FindSubmatch(content)
	if len(match) > index {
		return strings.TrimSpace(string(match[index]))
	}
	return ""
}
