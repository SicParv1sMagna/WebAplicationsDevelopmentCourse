package repository

import (
	"errors"
	"project/internal/model"

	"gorm.io/gorm"
)

func (r *Repository) CreateMarkdown(md model.Markdown) error {
	err := r.db.Table("Markdown").Create(&md).Error

	return err
}

func (r *Repository) GetAllMarkdowns(name string) ([]model.Markdown, error) {
	name = "%" + name + "%"
	var markdowns []model.Markdown

	err := r.db.Table("Markdown").Where(`("Status" = 'Активен' OR "Status" = 'Черновик') AND LOWER("Name") LIKE LOWER(?)`, name).Find(&markdowns).Error

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

	markdownContributor := model.MarkdownContributor{
		Markdown_ID:    id,
		Contributor_ID: uint(contributor.Contributor_ID),
		Status:         "Читатель",
	}

	err = r.db.Table("document_request").Create(&markdownContributor).Error
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

func (r *Repository) AddMdToLastReader(markdown_id, contributor_id uint) (model.MarkdownContributor, model.Markdown, error) {
	var markdown model.Markdown
	var markdownReader model.MarkdownContributor
	err := r.db.
		Table("document_request").
		Where(`"Status" = ?`, "Читатель").
		First(&markdownReader).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return markdownReader, markdown, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		markdownReader := model.MarkdownContributor{
			Markdown_ID:    markdown_id,
			Contributor_ID: contributor_id,
			Status:         "Читатель",
		}
		if err := r.db.Table("document_request").Create(&markdownReader).Error; err != nil {
			return markdownReader, markdown, err
		}
	}

	if err := r.db.Table("Markdown").First(&markdown, markdown_id).Error; err != nil {
		return markdownReader, markdown, err
	}

	if err := r.db.Table("document_request").Create(model.MarkdownContributor{
		Contributor_ID: uint(markdownReader.Contributor_ID),
		Markdown_ID:    uint(markdown.Markdown_ID),
		Status:         "Редактор",
	}).Error; err != nil {
		return markdownReader, markdown, err
	}

	return markdownReader, markdown, err
}

func (r *Repository) DeleteContributorFromMd(id, cid uint) error {
	err := r.db.Table("Document_Request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Set(`"Status" = ?`, "Удален").Error

	return err
}

func (r *Repository) AddMarkdownIcon(id int, imageBytes []byte, contentType string) error {
	// err := r.minio.RemoveServiceImage(id)
	// if err != nil {
	// 	return err
	// }

	imageURL, err := r.minio.UploadServiceImage(id, imageBytes, contentType)
	if err != nil {
		return err
	}

	err = r.db.Table("Markdown").Where("Markdown_ID = ?", id).Update("PhotoURL", imageURL).Error
	if err != nil {
		return errors.New("ошибка обновления url изображения в БД")
	}

	return nil
}
