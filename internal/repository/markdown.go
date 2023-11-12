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

func (r *Repository) RequestContribution(contributor model.Contributor, id uint) error {
	err := r.db.Table("contributor").Create(&contributor).Error
	if err != nil {
		return err
	}

	err = r.db.Exec(`INSERT INTO document_request (markdown_id, contributor_id, Status) VALUES($1, $2, $3)`, id, contributor.Contributor_ID, "Читатель").Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SearchMarkdown(query string) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	sqlname := "%" + query + "%"

	if err := r.db.Table("Markdown").Where(`"Name" LIKE ? AND "Status" = 'Активен'`, sqlname).Find(&markdowns).Error; err != nil {
		return nil, err
	}

	return markdowns, nil
}
