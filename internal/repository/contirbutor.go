package repository

import (
	"fmt"
	"project/internal/model"
	"time"
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

func (r *Repository) GetContributorsByMarkdownID(email string, status string, start_date string, end_date string, markdownID uint) ([]model.Contributor, error) {
	var contributors []model.Contributor

	query := r.db.
		Table("contributor c").
		Select("c.*").
		Joins("INNER JOIN document_request dr ON c.contributor_id = dr.contributor_id").
		Where(`dr.markdown_id = ? AND dr."Status" != 'Удален'`, markdownID)

	// Add filters based on the provided parameters
	if email != "" {
		query = query.Where("c.email = ?", email)
	}

	if status != "" {
		query = query.Where(`dr."Status" = ?`, status)
	}

	if start_date != "" {
		query = query.Where("c.created_date >= ?", start_date)
	}

	if end_date != "" {
		query = query.Where("c.created_date <= ?", end_date)
	}
	fmt.Println(query)
	// Execute the query and scan the result into the contributors slice
	if err := query.Find(&contributors).Error; err != nil {
		return nil, err
	}

	return contributors, nil
}

func (r *Repository) UpdateContributorAccess(id, cid uint, newStatus string) error {
	err := r.db.Table("Document_Request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Set(`"Status" = ?`, newStatus).Error
	if err != nil {
		return err
	}

	err = r.db.Table("contributor").Where("contributor_id = ?", cid).Set(`"formed_date" = ?`, time.Now()).Error
	if err != nil {
		return err
	}

	err = r.db.Table("contributor").Where("contributor_id = ?", cid).Set(`"completion_date" = ?`, time.Now()).Error

	return err
}

func (r *Repository) GetAllContributors(email, status, start_date, end_date string) ([]model.Contributor, error) {
	var contributors []model.Contributor

	// Start building the query
	query := r.db.Table("contributor").
		Select(`contributor.contributor_id, contributor.user_id, contributor.created_date, contributor.formed_date, contributor.completion_date, contributor.email, document_request."Status"`).
		Joins("LEFT JOIN document_request ON contributor.contributor_id = document_request.contributor_id").
		Order("contributor.created_date DESC")

	// Add filters based on the provided parameters
	if email != "" {
		query = query.Where("contributor.email = ?", email)
	}

	if status != "" {
		query = query.Where(`document_request."Status" = ?`, status)
	}

	if start_date != "" {
		query = query.Where("contributor.created_date >= ?", start_date)
	}

	if end_date != "" {
		query = query.Where("contributor.created_date <= ?", end_date)
	}

	// Execute the query and scan the result into the contributors slice
	if err := query.Find(&contributors).Error; err != nil {
		return nil, err
	}
	return contributors, nil
}

func (r *Repository) UpdateContributorData(id uint, email string) error {
	err := r.db.Table("contributor").Where("contributor_id = ?", id).Set("email = ?", email).Error

	return err
}
