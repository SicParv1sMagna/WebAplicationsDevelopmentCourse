package repository

import "project/internal/model"

func (r *Repository) GetMarkdownsByUserID(userID int) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	// SQL-запрос с явным SQL-кодом для выполнения JOIN через таблицы Document_Request и Contributor
	sqlQuery := `
        SELECT DISTINCT "Markdown".*
        FROM "Markdown"
        INNER JOIN Document_Request ON "Markdown".Markdown_ID = Document_Request.Markdown_ID
        INNER JOIN Contributor ON Document_Request.Contributor_ID = Contributor.Contributor_ID
        WHERE Contributor.User_ID = ?`

	// Выполняем запрос и сканируем результаты в структуру Markdown
	if err := r.db.Raw(sqlQuery, userID).Scan(&markdowns).Error; err != nil {
		return nil, err
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
