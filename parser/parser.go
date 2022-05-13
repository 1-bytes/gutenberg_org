package parser

import (
	"bbc_com/pkg/downloader"
	"regexp"
	"strconv"
)

var content string

type Parser struct {
	ID            string
	Author        string
	Title         string
	Language      string
	ReleaseDate   string
	DownloadCount int
	CoverImage    string
	response      string
}

var authorRe = regexp.MustCompile(`<a href="/ebooks/author[^>]*?">([^<]*?)</a></td>`)
var titleRe = regexp.MustCompile(`<meta name="title" content="([^"]*?)">`)
var languageRe = regexp.MustCompile(`<th>Language</th>\s<td>(<a href=[^>]*?>)?([^<]*?)(</a>)?</td>`)
var releaseDateRe = regexp.MustCompile(`<th>Release Date</th>\s<td[^>]*?>([^<]*?)</td>`)
var downloadCountRe = regexp.MustCompile(`<td itemprop="interactionCount">(\d*?) downloads in`)
var coverImageRe = regexp.MustCompile(`<img class="cover-art" src="([^"]*?)"\s`)

func (p *Parser) GetDetail(resp string) {
	p.response = resp
	p.Author = p.author()
	p.Title = p.title()
	p.Language = p.language()
	p.ReleaseDate = p.releaseDate()
	p.DownloadCount = p.downloadCount()
	p.CoverImage = p.coverImage()
}

// author 作者信息
func (p *Parser) author() string {
	match := authorRe.FindStringSubmatch(p.response)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// title 标题
func (p *Parser) title() string {
	match := titleRe.FindStringSubmatch(p.response)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// language 语言
func (p *Parser) language() string {
	match := languageRe.FindStringSubmatch(p.response)
	if len(match) > 0 {
		return match[2]
	}
	return ""
}

// releaseDate 发布日期
func (p *Parser) releaseDate() string {
	match := releaseDateRe.FindStringSubmatch(p.response)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// downloadCount 获取下载次数
func (p *Parser) downloadCount() int {
	match := downloadCountRe.FindStringSubmatch(p.response)
	if len(match) > 0 {
		c, err := strconv.Atoi(match[1])
		if err == nil {
			return c
		}
	}
	return 0
}

// releaseDate 发布日期
func (p *Parser) coverImage() string {
	coverMatch := coverImageRe.FindStringSubmatch(p.response)
	if len(coverMatch) > 0 {
		idRe := regexp.MustCompile(`/epub/([^/]*?)/`)
		idMatch := idRe.FindStringSubmatch(coverMatch[1])
		id := idMatch[1]
		downloader.DownloadImage(coverMatch[1], "files/cover_image/"+id+".jpg")
		return id + ".jpg"
	}
	return ""
}
