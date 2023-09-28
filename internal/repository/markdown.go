package repository

import "project/internal/model"

// Вытягиваем все Markdown - документы пользователя
func (r *Repository) GetAllNotes(userID uint) ([]model.Markdown, error) {
	// Инициализируем массив md
	var md []model.Markdown
	// Делаем запрос к базе данных на получение списка md
	err := r.db.Table("Markdown").Where("Moderator_ID = ?", userID).Find(&md).Error
	if err != nil {
		return nil, err
	}

	return md, nil
}

// Получаем документ по его ID
func (r *Repository) GetMarkdownById(mdID uint) (model.Markdown, error) {
	var md model.Markdown

	err := r.db.Table("Markdown").Where("Markdown_ID = ?", mdID).Find(&md).Error
	if err != nil {
		return md, err
	}

	return md, nil
}

func (r *Repository) DeleteMarkdownById(mdID uint) error {
	if err := r.db.Exec(`UPDATE "Markdown" SET "Status"='Удален' WHERE Markdown_ID = ?`, mdID).Error; err != nil {
		return err
	}

	return nil
}
