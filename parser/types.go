package parser

type ElasticSearchData struct {
	ID           string `json:"id"`
	SourceDomain string `json:"source_domain"`
	SourceUrl    string `json:"source_url"`
	Paragraph    struct {
		EN string `json:"en"`
	} `json:"paragraph"`
}
