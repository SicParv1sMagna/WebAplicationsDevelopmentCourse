package repository

import "project/internal/model"

func (r *Repository) GetContributorByID(id uint) (model.Contributor, error) {
	var contributor model.Contributor

	err := r.db.Table(`contributor`).Where("Contributor_ID = ?", id).First(&contributor).Error
	if err != nil {
		return contributor, err
	}

	return contributor, nil
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
