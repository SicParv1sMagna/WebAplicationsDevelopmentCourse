package usecase

import (
	"errors"
	"io"
	"mime/multipart"
	"project/internal/model"
	"project/internal/pkg/validators"
	"time"
)

func (uc *UseCase) CreateMarkdown(markdown model.Markdown, id uint) (model.Markdown, error) {
	if err := validators.ValidateMarkdown(markdown); err != nil {
		return model.Markdown{}, err
	}

	markdown.User_ID = id

	if err := uc.Repository.CreateMarkdown(markdown); err != nil {
		return model.Markdown{}, errors.New("ошибка при создании маркдауна")
	}

	return markdown, nil
}

func (uc *UseCase) GetAllMarkdown(name string) ([]model.Markdown, error) {
	md, err := uc.Repository.GetAllMarkdowns(name)
	if err != nil {
		return []model.Markdown{}, errors.New("ошибка при получении маркадунов")
	}

	return md, nil
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

func (uc *UseCase) AddMarkdownToContributor(mid, cid uint) error {
	if mid <= 0 || cid <= 0 {
		return errors.New("id не может быть отрицательным")
	}
	// черновик, требует подтверждения, в работе, отклонен, удален, завершен
	if err := uc.Repository.AddMarkdownToLastDraft(mid, cid); err != nil {
		return errors.New("ошибка при добавлении")
	}

	return nil
}

func (uc *UseCase) DeleteContributorFromMd(jsonData map[string]interface{}) error {
	cid, cidOk := jsonData["Contributor_ID"].(float64)
	mid, midOk := jsonData["Markdown_ID"].(float64)

	if !cidOk || !midOk {
		return errors.New("id мардауна или контрибьютора отсутствуют")
	}

	if cid <= 0 || mid <= 0 {
		return errors.New("id маркдауна или контрибьютора отрицательны")
	}

	err := uc.Repository.DeleteContributorFromMd(uint(mid), uint(cid))
	if err != nil {
		return errors.New("ошибка при удалении контрибьютора из маркдауна")
	}

	return nil
}

func (uc *UseCase) AddMarkdownIcon(id uint, image *multipart.FileHeader) error {
	if id <= 0 {
		return errors.New("id не может быть отрицательным")
	}

	file, err := image.Open()
	if err != nil {
		return errors.New("ошибка при открытии изображения")
	}

	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return errors.New("ошибка чтения изображения")
	}

	contentType := image.Header.Get("Content-Type")

	if err = uc.Repository.AddMarkdownIcon(int(id), imageBytes, contentType); err != nil {
		return errors.New("ошибка при добавлении иконки")
	}

	return nil
}

func (uc *UseCase) RequestContribution(uid, mid uint) error {
	if uid <= 0 || mid <= 0 {
		return errors.New("id пользователя или маркдауна не могут быть отрицательными")
	}

	candidates, err := uc.Repository.GetContributorsByMarkdownID("", "", "", "", uint(mid))
	if err != nil {
		return errors.New("ошибка при создании запроса на редактирование")
	}

	for i := 0; i < len(candidates); i++ {
		if uint(candidates[i].User_ID) == uid {
			return errors.New("вы уже подали запрос на редактирование")
		}
	}

	user, err := uc.Repository.GetUserById(uid)
	if err != nil {
		return errors.New("ошибка при создании запроса на редактирование")
	}

	contributor := model.Contributor{
		User_ID:      int(uid),
		Created_Date: time.Now(),
		Email:        user.Email,
	}

	err = uc.Repository.RequestContribution(contributor, uint(mid))
	if err != nil {
		return errors.New("ошибка при создании запроса на редактирование")
	}

	return nil
}
