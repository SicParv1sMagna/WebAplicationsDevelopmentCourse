package repository

import (
	"project/internal/model"
)

func (r *Repository) CreateMarkdown(md model.Markdown) error {
	err := r.db.Table("Markdown").Create(&md).Error

	return err
}

func (r *Repository) GetAllMarkdowns() ([]model.Markdown, error) {
	var markdowns []model.Markdown

	err := r.db.Table("Markdown").Find(&markdowns).Error

	return markdowns, err
}

func (r *Repository) GetMarkdownById(mdID uint) (model.Markdown, error) {
	var markdown model.Markdown

	err := r.db.Table("Markdown").Where("Markdown_ID = ?", mdID).First(&markdown).Error
	if err != nil {
		return markdown, err
	}

	return markdown, nil
}

func (r *Repository) DeleteMarkdownById(mdID uint) error {
	err := r.db.Table("Markdown").Where("Markdown_ID = ?", mdID).Update("Status", "Удален").Error

	return err
}

func (r *Repository) UpdateMarkdownById(md model.Markdown) error {
	err := r.db.Table("Markdown").Where("Markdown_ID = ?", md.Markdown_ID).Updates(&md).Error

	return err
}
