package repository

import "project/internal/model"

func (r *Repository) GetMarkdownsByUserID(userID int) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	// SQL-запрос с явным SQL-кодом для выполнения JOIN через таблицы Document_Request и Contributor

	// Выполняем запрос и сканируем результаты в структуру Markdown
	if err := r.db.Table("Markdown").Where(`"Status" = 'Активен'`).Find(&markdowns).Error; err != nil {
		return markdowns, err
	}

	return markdowns, nil
}

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

func (r *Repository) SearchMarkdowns(userID int, name string) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	sqlname := "%" + name + "%"
	// Выполняем запрос с использованием GORM
	if err := r.db.Table("Markdown").Where(`"Name" LIKE ? AND "Status" = 'Активен'`, sqlname).Find(&markdowns).Error; err != nil {
		return nil, err
	}

	return markdowns, nil
}
