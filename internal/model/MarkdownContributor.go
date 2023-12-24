package model

type MarkdownContributor struct {
	Markdown_ID    uint `gorm: "primaryKey;column:markdown_id"`
	Contributor_ID uint `gorm:"primaryKey;column:contributor_id"`
}
