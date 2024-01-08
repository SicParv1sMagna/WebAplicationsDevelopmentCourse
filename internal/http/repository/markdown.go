package repository

import (
	"errors"
	"project/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (r *Repository) CreateMarkdown(md model.Markdown) (uint, error) {
	err := r.db.Exec(`INSERT INTO "Markdown" ("Name", "Content", "Status", "created_time", "user_id", "PhotoURL") VALUES ($1, $2, $3, $4, $5, $6)`, md.Name, md.Content, md.Status, md.Created_Time, md.User_ID, md.PhotoURL).Error

	return uint(md.Markdown_ID), err
}

func (r *Repository) GetAllMarkdowns(name string, userID int) ([]model.Markdown, uint, error) {
	var contributorId uint
	var contributor model.Contributor
	if userID != 0 {
		err := r.db.Table("contributor").Where(`user_id = ? and "Status" = ?`, userID, "Черновик").First(&contributor).Error
		if err != nil {
			contributorId = 0
		} else {
			contributorId = uint(contributor.Contributor_ID)
		}
	} else {
		contributorId = 0
	}

	name = "%" + name + "%"
	var markdowns []model.Markdown

	err := r.db.Table("Markdown").Where(`("Status" = 'Активен' OR "Status" = 'Черновик') AND LOWER("Name") LIKE LOWER(?)`, name).Order("markdown_id DESC").Find(&markdowns).Error

	return markdowns, contributorId, err
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

func (r *Repository) AddMarkdownToLastDraft(markdownID, userID uint) error {
	var user model.User
	if err := r.db.Table("User").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	// Находим последнюю миссию с mission_status = "Draft"
	var lastDraftContributor model.Contributor
	dbErr := r.db.
		Table("contributor").
		Order("created_date DESC").
		Where(`user_id = ? AND "Status" = ?`, userID, "Черновик").
		First(&lastDraftContributor).
		Error

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return dbErr
	}

	// Если миссии с mission_status = "Draft" нет, создаем новую
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		currentTime := time.Now()

		lastDraftContributor = model.Contributor{
			User_ID:         int(userID),
			Created_Date:    &currentTime,
			Formed_Date:     nil,
			Completion_Date: nil,
			Status:          "Черновик",
			Email:           user.Email,
		}
		if err := r.db.Table("contributor").Create(&lastDraftContributor).Error; err != nil {
			return err
		}
	}

	// Получаем образец из базы данных по его идентификатору
	var newMarkdown model.Markdown
	if err := r.db.Table("Markdown").First(&newMarkdown, markdownID).Error; err != nil {
		return err
	}

	// Добавляем образец в миссию
	if err := r.db.Table("document_request").Create(&model.MarkdownContributor{
		Markdown_ID:    uint(newMarkdown.Markdown_ID),
		Contributor_ID: uint(lastDraftContributor.Contributor_ID),
	}).Error; err != nil {
		// Проверяем, является ли ошибка уникальным ключом
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			// Здесь обрабатываем случай дубликата ключа, если это произошло
			return errors.New("Образец уже добавлен в миссию")
		}
		return err
	}

	return nil
}

func (r *Repository) DeleteContributorFromMd(id, uid uint) error {
	var lastDraftContributor model.Contributor
	dbErr := r.db.
		Table("contributor").
		Order("formed_date desc").
		Where(`"Status" = ? AND user_id = ?`, "Черновик", uid).
		First(&lastDraftContributor).
		Error

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return dbErr
	}

	// Если миссии с mission_status = "Draft" нет, возвращаем ошибку
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return errors.New("Контрибьютор со статусом черновик не найдена")
	}

	// Удаляем образец из миссии
	if err := r.db.Exec("DELETE FROM document_request WHERE contributor_id = ? AND markdown_id = ?", lastDraftContributor.Contributor_ID, id).Error; err != nil {
		return err
	}

	return nil
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
	err := r.db.Table("contributor").Where("user_id = ?", user_id).Update("formed_date", time.Now()).Update("Status", "Требует подтверждения").Error
	if err != nil {
		return err
	}

	return nil
}
