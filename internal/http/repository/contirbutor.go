package repository

import (
	"project/internal/model"
	"time"
)

func (r *Repository) GetContributorByUserID(userID uint) (model.Contributor, error) {
	var contributor model.Contributor

	err := r.db.Table("contributor").Where("user_id = ?", userID).First(&contributor).Error
	if err != nil {
		return model.Contributor{}, err
	}

	return contributor, nil
}

func (r *Repository) GetContributorByID(id uint) (model.Contributor, []model.Markdown, error) {
	var contributor model.Contributor
	var markdowns []model.Markdown

	err := r.db.
		Joins("JOIN document_request ON contributor.contributor_id = document_request.contributor_id").
		Joins(`JOIN "Markdown" ON document_request.markdown_id = "Markdown".markdown_id`).
		Where("contributor.contributor_id = ?", id).
		Table("contributor").
		Select(`contributor.*, "Markdown".*`).
		Find(&markdowns).Error
	if err != nil {
		return model.Contributor{}, nil, err
	}

	return contributor, markdowns, nil
}

func (r *Repository) GetContributorsByMarkdownID(email string, status string, start_date string, end_date string, markdownID uint) ([]model.Contributor, error) {
	var contributors []model.Contributor

	query := r.db.
		Table("contributor c").
		Select("c.*").
		Joins("INNER JOIN document_request dr ON c.contributor_id = dr.contributor_id").
		Where(`dr.markdown_id = ? AND c."Status" != 'Удален'`, markdownID)

	// Add filters based on the provided parameters
	if email != "" {
		query = query.Where("c.email = ?", email)
	}

	if status != "" {
		query = query.Where(`c."Status" = ?`, status)
	}

	if start_date != "" {
		query = query.Where("c.created_date >= ?", start_date)
	}

	if end_date != "" {
		query = query.Where("c.created_date <= ?", end_date)
	}

	// Execute the query and scan the result into the contributors slice
	if err := query.Find(&contributors).Error; err != nil {
		return nil, err
	}

	return contributors, nil
}

func (r *Repository) UpdateContributorAccessByModerator(cid, moderatorId uint, newStatus string) error {
	currentTime := time.Now()

	var user model.User

	err := r.db.Table("User").Where("user_id = ?", moderatorId).First(&user).Error
	if err != nil {
		return err
	}

	err = r.db.Table("contributor").Where("contributor_id = ?", cid).Update(`"completion_date"`, currentTime).Update(`"Status"`, newStatus).Update("approved_by", user.Email).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateContributorAccessByAdmin(id, cid uint, newStatus string) error {
	err := r.db.Table("document_request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Update(`Status`, newStatus).Error
	if err != nil {
		return err
	}

	currentTime := time.Now()

	err = r.db.Table("contributor").Where("contributor_id = ?", cid).Update(`"formed_date"`, currentTime).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllContributors(email, status, start_date, end_date string) ([]model.Contributor, error) {
	var contributors []model.Contributor

	query := r.db.Table("contributor")
	// Add filters based on the provided parameters
	if email != "" {
		email = "%" + email + "%"
		query = query.Where("LOWER(contributor.email) LIKE LOWER(?)", email)
	}

	if status != "" {
		query = query.Where(`contributor."Status" = ?`, status)
	}

	if start_date != "" {
		query = query.Where("contributor.created_date >= ?", start_date)
	}

	if end_date != "" {
		query = query.Where("contributor.created_date <= ?", end_date)
	}

	// Execute the query and scan the result into the contributors slice
	if err := query.Order("contributor_id DESC").Where(`contributor."Status" != ? AND contributor."Status" != ?`, "Удален", "Черновик").Find(&contributors).Error; err != nil {
		return nil, err
	}

	return contributors, nil
}

func (r *Repository) UpdateContributorData(id uint, email string) error {
	err := r.db.Table("contributor").Where("contributor_id = ?", id).Set("email = ?", email).Error

	return err
}
