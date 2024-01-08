package usecase

import (
	"errors"
	"io"
	"mime/multipart"
	"project/internal/model"
	"project/internal/pkg/validators"
)

func (uc *UseCase) CreateMarkdown(markdown model.Markdown, id uint) (model.Markdown, error) {
	if err := validators.ValidateMarkdown(markdown); err != nil {
		return model.Markdown{}, err
	}

	markdown.User_ID = id

	id, err := uc.Repository.CreateMarkdown(markdown)
	if err != nil {
		return model.Markdown{}, errors.New("ошибка при создании маркдауна")
	}

	markdown.Markdown_ID = int(id)

	return markdown, nil
}

func (uc *UseCase) GetAllMarkdown(name string, userID int) ([]model.Markdown, uint, error) {
	md, id, err := uc.Repository.GetAllMarkdowns(name, userID)
	if err != nil {
		return []model.Markdown{}, id, errors.New("ошибка при получении маркадунов")
	}

	return md, id, nil
}

func (uc *UseCase) GetMarkdown(id int) (model.Markdown, error) {
	if id <= 0 {
		return model.Markdown{}, errors.New("id должен быть не более 0")
	}

	markdown, err := uc.Repository.GetMarkdownById(uint(id))
	if err != nil {
		return model.Markdown{}, errors.New("ошибка при получении маркдауна")
	}

	return markdown, nil
}

func (uc *UseCase) DeleteMarkdown(id int) error {
	if id <= 0 {
		return errors.New("id должен быть не более 0")
	}

	err := uc.Repository.DeleteMarkdownById(uint(id))
	if err != nil {
		return errors.New("ошибка при удалении маркдауна")
	}

	return nil
}

func (uc *UseCase) UpdateMarkdown(markdown map[string]interface{}) error {
	mdID, ok := markdown["Markdown_ID"].(float64)
	if !ok || mdID <= 0 {
		return errors.New("id маркадуна отсутствует или меньше 1")
	}
	name, _ := markdown["Name"].(string)
	content, _ := markdown["Content"].(string)

	candidate, err := uc.Repository.GetMarkdownById(uint(mdID))
	if err != nil {
		return errors.New("ошибка получения маркдауна")
	}

	if len(name) >= 0 {
		candidate.Name = name
	}

	if len(content) >= 0 {
		candidate.Content = content
	}

	if err = uc.Repository.UpdateMarkdownById(candidate); err != nil {
		return errors.New("ошибка при обнволении маркдауна")
	}

	return nil
}

func (uc *UseCase) AddMarkdownToContributor(mid, uid uint) error {
	if mid <= 0 || uid <= 0 {
		return errors.New("id не может быть отрицательным")
	}
	// черновик, требует подтверждения, в работе, отклонен, удален, завершен
	if err := uc.Repository.AddMarkdownToLastDraft(mid, uid); err != nil {
		return errors.New("ошибка при добавлении")
	}

	return nil
}

func (uc *UseCase) DeleteContributorFromMd(markdownID, userID int) error {
	err := uc.Repository.DeleteContributorFromMd(uint(markdownID), uint(userID))
	if err != nil {
		return errors.New("ошибка при удалении контрибьютора из маркдауна")
	}

	return nil
}

func (uc *UseCase) AddMarkdownIcon(id uint, image *multipart.FileHeader) (string, error) {
	if id <= 0 {
		return "", errors.New("id не может быть отрицательным")
	}

	file, err := image.Open()
	if err != nil {
		return "", errors.New("ошибка при открытии изображения")
	}

	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("ошибка чтения изображения")
	}

	contentType := image.Header.Get("Content-Type")

	url, err := uc.Repository.AddMarkdownIcon(int(id), imageBytes, contentType)
	if err != nil {
		return "", errors.New("ошибка при добавлении иконки")
	}

	return url, nil
}

func (uc *UseCase) RequestContribution(uid, mid uint) error {
	if uid <= 0 || mid <= 0 {
		return errors.New("id пользователя или маркдауна не могут быть отрицательными")
	}

	err := uc.Repository.RequestContribution(uint(uid), uint(mid))
	if err != nil {
		return errors.New("ошибка при создании запроса на редактирование")
	}

	return nil

}
