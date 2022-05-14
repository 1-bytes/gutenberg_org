package model

import "spider/bootstrap"

// 段落内容
type DictArticleParagraph struct {
	ID            int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ElasticID     string `gorm:"column:elastic_id;default:0;NOT NULL"`      // ElasticSearch ID
	ArticleID     int    `gorm:"column:article_id;default:0;NOT NULL"`      // 文章ID
	ByteCount     int    `gorm:"column:byte_count;default:0;NOT NULL"`      // 段落长度
	Fre           string `gorm:"column:fre;NOT NULL"`                       // Flesch Reading Ease.FRE数值越高，文章就越简单，可读性也越高。
	Fkgl          string `gorm:"column:fkgl;NOT NULL"`                      // Flesch–Kincaid Grade Level.FKGL数值越高，文章就越复杂，文章的可读性也就越低。
	SchoolLvClass int    `gorm:"column:school_lv_class;default:0;NOT NULL"` // 依据fre给出的评级
	SchoolLvlName string `gorm:"column:school_lvl_name;NOT NULL"`           // slv对应的学校等级
	TechWordLv    string `gorm:"column:tech_word_lv;NOT NULL"`              // 人教评级
	CefrWordLv    string `gorm:"column:cefr_word_lv;NOT NULL"`              // cefr评级
	Status        int    `gorm:"column:status;default:1"`                   // 状态：1-正常；0-删除
}

// TableName 获取表名
func (m *DictArticleParagraph) TableName() string {
	return "dict_article_paragraph_test"
}

// BatchCreate 批量插入数据
func (m *DictArticleParagraph) BatchCreate(models []DictArticleParagraph) error {
	tx := bootstrap.DB.Table(m.TableName()).Create(models)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
