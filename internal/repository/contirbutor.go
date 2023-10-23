package repository

import (
	"errors"
	"fmt"
	"project/internal/model"

	"gorm.io/gorm"
)

func (r *Repository) GetContributorByID(id uint) (model.Contributor, []model.Markdown, error) {
	var contributor model.Contributor
	var markdowns []model.Markdown

	err := r.db.Table(`contributor`).Where("Contributor_ID = ?", id).First(&contributor).Error
	if err != nil {
		return contributor, markdowns, err
	}

	err = r.db.
		Joins(`JOIN document_request ON contributor.Contributor_ID = document_request.Contributor_ID`).
		Joins(`JOIN "Markdown" ON document_request.Markdown_ID = "Markdown".Markdown_ID`).
		Where("contributor.Contributor_ID = ?", id).
		Table("contributor"). // Add this line to specify the table name
		Select(`contributor.*, "Markdown".*`).
		Find(&markdowns).Error

	if err != nil {
		return contributor, markdowns, err
	}

	return contributor, markdowns, nil
}

func (r *Repository) GetContributorsByMarkdownID(markdownID uint) ([]model.Contributor, error) {
	var contributors []model.Contributor

	err := r.db.
		Table("contributor c").
		Select("c.*").
		Joins("INNER JOIN Document_Request dr ON c.Contributor_ID = dr.Contributor_ID").
		Where(`dr.Markdown_ID = ? AND dr."Status" != 'Удален'`, markdownID).
		Find(&contributors).Error

	if err != nil {
		return nil, err
	}

	return contributors, nil
}

func (r *Repository) DeleteContributorFromMd(id, cid uint) error {
	err := r.db.Table("Document_Request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Set(`"Status" = ?`, "Удален").Error

	return err
}

func (r *Repository) UpdateContributorAccess(id, cid uint, newStatus string) error {
	err := r.db.Table("Document_Request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Set(`"Status" = ?`, newStatus).Error

	return err
}

func (r *Repository) AddMdToLastReader(id uint) (model.MarkdownContributor, model.Markdown, error) {
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
			Markdown_ID:    id,
			Contributor_ID: 1,
			Status:         "Читатель",
		}
		if err := r.db.Table("document_request").Create(&markdownReader).Error; err != nil {
			return markdownReader, markdown, err
		}
	}

	if err := r.db.Table("Markdown").First(&markdown, id).Error; err != nil {
		return markdownReader, markdown, err
	}
	fmt.Println(markdown)
	fmt.Println(markdownReader)

	if err := r.db.Table("document_request").Create(model.MarkdownContributor{
		Contributor_ID: uint(markdownReader.Contributor_ID),
		Markdown_ID:    uint(markdown.Markdown_ID),
		Status:         "Редактор",
	}).Error; err != nil {
		return markdownReader, markdown, err
	}

	return markdownReader, markdown, err
}
