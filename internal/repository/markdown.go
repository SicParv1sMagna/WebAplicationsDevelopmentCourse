package repository

import (
	"errors"
	"fmt"
	"project/internal/model"
)

func (r *Repository) CreateMarkdown(md model.Markdown, userID uint) error {
	var markdownID int
	sql := `INSERT INTO "Markdown" ("Name") VALUES (?) RETURNING Markdown_ID`

	err := r.db.Raw(sql, md.Name).Scan(&markdownID).Error
	if err != nil {
		return err
	}

	// Попытка вставить запись в Contributor
	sql = `INSERT INTO Contributor (Moderator_ID, User_ID) VALUES (?, ?) ON CONFLICT (Moderator_ID, User_ID) DO NOTHING`
	if err := r.db.Exec(sql, userID, userID).Error; err != nil {
		return err
	}

	// Запрос к Contributor для получения ID записи
	var contributorID int
	sql = `SELECT Contributor_ID FROM Contributor WHERE Moderator_ID = ? AND User_ID = ?`
	if err := r.db.Raw(sql, userID, userID).Scan(&contributorID).Error; err != nil {
		return err
	}

	// Вставка записи в Document_Request
	sql = `INSERT INTO Document_Request (Markdown_ID, Contributor_ID) VALUES (?, ?)`
	if err := r.db.Exec(sql, markdownID, contributorID).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllMarkdowns(userID uint) ([]model.Markdown, error) {
	var markdowns []model.Markdown

	// Выполняем SQL-запрос, объединяя таблицы Document_Request и Markdown
	// и выбираем все записи, где Contributor_ID равно userID
	err := r.db.
		Table("Document_Request dr").
		Joins(`INNER JOIN "Markdown" m ON dr.Markdown_ID = m.Markdown_ID`).
		Where("dr.Contributor_ID = ?", userID).
		Select("m.*").
		Find(&markdowns).Error

	if err != nil {
		return nil, err
	}

	return markdowns, nil
}

func (r *Repository) GetMarkdownById(mdID, userID uint) (model.Markdown, error) {
	var md model.Markdown

	// Используйте метод Raw для выполнения SQL-запроса с JOIN
	sql := `
        SELECT m.*
        FROM "Markdown" m
        JOIN Document_Request dr ON m.markdown_id = dr.markdown_id
        JOIN Contributor c ON dr.contributor_id = c.contributor_id
        WHERE m.markdown_id = ? AND c.user_id = ?
    `

	// Выполняем запрос и сканируем результат в структуру Markdown
	if err := r.db.Raw(sql, mdID, userID).Scan(&md).Error; err != nil {
		return md, err
	}

	// Если mdID не существует или не принадлежит пользователю userID, вернуть ошибку
	if md.Markdown_ID == 0 {
		return md, fmt.Errorf("markdown не найден")
	}

	return md, nil
}

func (r *Repository) DeleteMarkdownById(mdID, userID uint) error {
	// Проверяем, что существует запись Contributor, связанная с данным маркдауном и пользователем
	var isAuthorized bool
	if err := r.db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM contributor c
			JOIN document_request dr ON c.contributor_id = dr.contributor_id
			WHERE c.user_id = ? AND dr.markdown_id = ?
		)
	`, userID, mdID).Row().Scan(&isAuthorized); err != nil {
		return err
	}

	if !isAuthorized {
		return errors.New("пользователь не авторизован для удаления этого маркдауна")
	}

	// Выполняем SQL-запрос для обновления статуса Markdown на "Удален"
	sql := `UPDATE "Markdown" SET "Status" = 'Удален' WHERE markdown_id = ?`

	if err := r.db.Exec(sql, mdID).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateMarkdownById(md model.Markdown, userID uint) error {
	// Выполните SQL-запрос для обновления маркдауна с проверкой через JOIN
	sql := `
        UPDATE "Markdown" m
        SET "Name" = ?, "Content" = ?
        FROM Document_Request dr
        JOIN Contributor c ON dr.contributor_id = c.contributor_id
        WHERE m.markdown_id = ? AND dr.markdown_id = m.markdown_id
          AND c.user_id = ?
    `

	// Выполните запрос на обновление маркдауна
	if err := r.db.Exec(sql, md.Name, md.Content, md.Markdown_ID, userID).Error; err != nil {
		return err
	}

	// Проверьте, обновлен ли хотя бы один маркдаун
	if r.db.RowsAffected == 0 {
		return fmt.Errorf("markdown не найден или не принадлежит пользователю")
	}

	return nil
}
