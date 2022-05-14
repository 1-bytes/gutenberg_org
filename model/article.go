package model

// 文章表，内容对应在paragraph
type DictArticle struct {
	ID                  uint   `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Type                string `gorm:"column:type;default:0;NOT NULL"`           // 类型 0:文章 1:图书
	Title               string `gorm:"column:title;NOT NULL"`                    // 标题
	Author              string `gorm:"column:author;NOT NULL"`                   // 作者
	CoverImage          string `gorm:"column:cover_image;NOT NULL"`              // 封面图
	Category            string `gorm:"column:category;default:0;NOT NULL"`       // 分类 0:未知 1:新闻 2:科普 3:笑话 4:小说 5:娱乐 6:诗歌 7:散文 8:故事 9:演讲 10.户外运动 11.上班族 12.旅游 13.未来 14.气象 15.文化
	DownloadCount       int    `gorm:"column:download_count;default:0;NOT NULL"` // 下载数量
	ReleaseDate         string `gorm:"column:release_date;NOT NULL"`             // 发布时间
	MostRecentlyUpdated string `gorm:"column:most_recently_updated;NOT NULL"`    // 最近更新时间
	SourceDomain        string `gorm:"column:source_domain;default:0;NOT NULL"`  // 来源域名 (对应 dict_source_domain_map 表)
	Status              string `gorm:"column:status;default:1;NOT NULL"`         // 状态：1-正常；0-删除
	Rank                int    `gorm:"column:rank;default:99;NOT NULL"`          // 排名：分数越小，排名越高
}

func (m *DictArticle) TableName() string {
	return "dict_article_test"
}
