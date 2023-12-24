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

func (r *Repository) GetContributorByID(id uint, status, start_date, end_date string) (model.Contributor, []model.Markdown, error) {
	var contributor model.Contributor
	var markdowns []model.Markdown

	// Query to fetch contributor details
	err := r.db.Table(`contributor`).Where("Contributor_ID = ?", id).First(&contributor).Error
	if err != nil {
		return contributor, markdowns, err
	}

	// Query to fetch markdowns based on provided parameters
	query := r.db.
		Joins(`JOIN document_request ON contributor.Contributor_ID = document_request.Contributor_ID`).
		Joins(`JOIN "Markdown" ON document_request.Markdown_ID = "Markdown".Markdown_ID`).
		Where("contributor.Contributor_ID = ?", id).
		Table("contributor").
		Select(`contributor.*, "Markdown".*, document_request."Status"`)

	// Add conditions based on parameters
	if status != "" {
		query = query.Where(`document_request."Status" = ?`, status)
	}

	if start_date != "" {
		query = query.Where("contributor.created_date >= ?", start_date)
	}

	if end_date != "" {
		query = query.Where("contributor.created_date <= ?", end_date)
	}

	// Execute the query and scan the result into the markdowns slice
	err = query.Find(&markdowns).Error
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
	// Execute the query and scan the result into the contributors slice
	if err := query.Find(&contributors).Error; err != nil {
		return nil, err
	}

	return contributors, nil
}

func (r *Repository) UpdateContributorAccessByModerator(id, cid uint, newStatus string) error {
	err := r.db.Table("document_request").Where("Contributor_ID = ? AND Markdown_ID = ?", cid, id).Update(`Status`, newStatus).Error
	if err != nil {
		return err
	}

	currentTime := time.Now()

	err = r.db.Table("contributor").Where("contributor_id = ?", cid).Update(`"completion_date"`, currentTime).Error
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

func (r *Repository) GetAllContributors(email, status, start_date, end_date string) ([]model.ContributorWithStatus, error) {
	var contributors []model.ContributorWithStatus

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

	for i := range contributors {
		// Assuming each contributor has an array of markdowns
		for j := range contributors[i].Markdowns {
			if contributors[i].Created_Date != nil {
				contributors[i].Markdowns[j].Created_Date = contributors[i].Created_Date
			}
			if contributors[i].Formed_Date != nil {
				contributors[i].Markdowns[j].Formed_Date = contributors[i].Formed_Date
			}
			if contributors[i].Completion_Date != nil {
				contributors[i].Markdowns[j].Completion_Date = contributors[i].Completion_Date
			}
		}
	}

	return contributors, nil
}

func (r *Repository) UpdateContributorData(id uint, email string) error {
	err := r.db.Table("contributor").Where("contributor_id = ?", id).Set("email = ?", email).Error

	return err
}

func (r *Repository) GetContributorStatus(cid, mid uint) (string, error) {
	var markdownContributor model.MarkdownContributor
	err := r.db.Table("document_request").Where("markdown_id = ? AND contributor_id = ?", mid, cid).First(&markdownContributor).Error
	if err != nil {
		return "", err
	}

	return markdownContributor.Status, nil
}
