package model

type MarkdownContributor struct {
	Markdown_ID    uint
	Contributor_ID uint
	Status         string `gorm:"column:Status"`
}
