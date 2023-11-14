package repository

import (
	"project/internal/model"
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

func (r *Repository) UpdateContributorAccess(id, cid uint, newStatus string) error {
	err := r.db.Table("Document_Request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Set(`"Status" = ?`, newStatus).Error

	return err
}
