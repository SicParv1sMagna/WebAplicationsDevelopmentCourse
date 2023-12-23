package repository

import (
	"errors"
	"project/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) CreateMarkdown(md model.Markdown) (uint, error) {
	err := r.db.Table("Markdown").Create(&md).Error

	return uint(md.Markdown_ID), err
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

func (r *Repository) SearchMarkdown(query string) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	sqlname := "%" + query + "%"

	if err := r.db.Table("Markdown").Where(`"Name" LIKE ? AND "Status" = 'Активен'`, sqlname).Find(&markdowns).Error; err != nil {
		return nil, err
	}

	return markdowns, nil
}

func isDuplicateKeyError(err error) bool {
	// Check the specific error codes that indicate a duplicate key violation
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint")
}

func (r *Repository) AddMarkdownToLastDraft(markdownID, userID uint) error {
	var contributor model.Contributor
	// Check if contributor exists
	err := r.db.Table("contributor").Where("user_id = ?", userID).First(&contributor).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// If contributor doesn't exist, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err := r.GetUserById(userID)
		if err != nil {
			return err
		}

		createdDate := time.Now()

		contributor = model.Contributor{
			User_ID:         int(userID),
			Email:           user.Email,
			Formed_Date:     nil,
			Created_Date:    &createdDate,
			Completion_Date: nil,
		}

		// Create new contributor
		err = r.db.Table("contributor").Create(&contributor).Error
		if err != nil {
			return err
		}

		err = r.db.Table("contributor").Where("user_id = ?", userID).First(&contributor).Error
		if err != nil {
			return err
		}
	}

	markdownContributor := model.MarkdownContributor{
		Markdown_ID:    markdownID,
		Contributor_ID: uint(contributor.Contributor_ID),
		Status:         "Черновик",
	}

	// Create Markdown contributor entry
	err = r.db.Table("document_request").Create(&markdownContributor).Error

	if err != nil {
		if isDuplicateKeyError(err) {
			err := r.db.Table("document_request").
				Where("contributor_id = ? AND markdown_id = ?",
					markdownContributor.Contributor_ID,
					markdownContributor.Markdown_ID).
				Updates(map[string]interface{}{"Status": "Черновик"}).Error

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteContributorFromMd(id, cid uint) error {
	err := r.db.Table("document_request").
		Where("contributor_id = ? AND markdown_id = ?", cid, id).
		Updates(map[string]interface{}{"Status": "Удален"}).Error

	return err
}

func (r *Repository) AddMarkdownIcon(id int, imageBytes []byte, contentType string) (string, error) {
	// err := r.minio.RemoveServiceImage(id)
	// if err != nil {
	// 	return err
	// }

	imageURL, err := r.minio.UploadServiceImage(id, imageBytes, contentType)
	if err != nil {
		return "", err
	}

	err = r.db.Table("Markdown").Where("Markdown_ID = ?", id).Update("PhotoURL", imageURL).Error
	if err != nil {
		return "", errors.New("ошибка обновления url изображения в БД")
	}

	return imageURL, nil
}

func (r *Repository) RequestContribution(user_id uint, markdown_id uint) error {
	var contributor model.Contributor

	err := r.db.Table("contributor").Where("user_id = ?", user_id).First(&contributor).Error
	if err != nil {
		return err
	}

	err = r.db.Table("document_request").Where("markdown_id = ? AND contributor_id = ?", markdown_id, contributor.Contributor_ID).Update("Status", "Требует подтверждения").Error
	if err != nil {
		return err
	}

	err = r.db.Table("contributor").Where("contributor_id = ?", contributor.Contributor_ID).Update("formed_date", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}
