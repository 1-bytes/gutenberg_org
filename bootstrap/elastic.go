package bootstrap

import (
	"bbc_com/pkg/config"
	pkgelastic "bbc_com/pkg/elastic"
	"github.com/olivere/elastic/v7"
	"time"
)

// SetupElastic 初始化 Elastic.
func SetupElastic() {
	pkgelastic.Options = []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetString("elastic.host")),
		elastic.SetBasicAuth(
			config.GetString("elastic.username"),
			config.GetString("elastic.password"),
		),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5 * time.Second),
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	}
}
