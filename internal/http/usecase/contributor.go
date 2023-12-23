package usecase

import (
	"errors"
	"project/internal/model"
)

func (uc *UseCase) GetContributor(id uint, start_date, end_date, status string) (model.Contributor, []model.Markdown, error) {
	if id <= 0 {
		return model.Contributor{}, []model.Markdown{}, errors.New("id не может быть отрицательным")
	}

	contributor, markdowns, err := uc.Repository.GetContributorByID(uint(id), status, start_date, end_date)
	if err != nil {
		return model.Contributor{}, []model.Markdown{}, errors.New("ошибка при получении данных")
	}

	return contributor, markdowns, nil
}

func (uc *UseCase) GetAllContributorsFromMarkdown(email, status, start_date, end_date string, id uint) ([]model.Contributor, error) {
	if id <= 0 {
		return []model.Contributor{}, errors.New("id не может быть отрицательным")
	}

	contributors, err := uc.Repository.GetContributorsByMarkdownID(email, status, start_date, end_date, id)
	if err != nil {
		return []model.Contributor{}, errors.New("ошибка при получении данных")
	}

	return contributors, nil
}

func (uc *UseCase) GetAllContributors(email, status, start_date, end_date string) ([]model.Contributor, error) {
	contributors, err := uc.Repository.GetAllContributors(email, status, start_date, end_date)
	if err != nil {
		return []model.Contributor{}, errors.New("ошибка при получении данных")
	}

	return contributors, nil
}

func (uc *UseCase) UpdateContributorAccessByModerator(jsonData map[string]interface{}) error {
	cid, ok := jsonData["Contributor_ID"].(float64)
	if !ok {
		return errors.New("поле контрибьютор должно быть передано")
	}
	if cid <= 0 {
		return errors.New("id контрибьютора не может быть отрицательным")
	}
	mid, ok := jsonData["Markdown_ID"].(float64)
	if !ok {
		return errors.New("поле маркдаун должно быть передано")
	}
	if mid <= 0 {
		return errors.New("id маркдауна не может быть отрицательным")
	}
	access, ok := jsonData["Access"].(string)
	if !ok {
		return errors.New("поле доступ должно быть заполнено")
	}
	if access == "Черновик" {
		return errors.New("нельзя вернуть статус в черновик")
	}
	if access == "В работе" || access == "Удален" {
		err := uc.Repository.UpdateContributorAccessByModerator(uint(mid), uint(cid), access)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *UseCase) UpdateContributorAccessByAdmin(jsonData map[string]interface{}) error {
	cid, ok := jsonData["Contributor_ID"].(float64)
	if !ok {
		return errors.New("поле контрибьютор должно быть передано")
	}
	if cid <= 0 {
		return errors.New("id контрибьютора не может быть отрицательным")
	}
	mid, ok := jsonData["Markdown_ID"].(float64)
	if !ok {
		return errors.New("поле маркдаун должно быть передано")
	}
	if mid <= 0 {
		return errors.New("id маркдауна не может быть отрицательным")
	}
	access, ok := jsonData["Access"].(string)
	if !ok {
		return errors.New("поле доступ должно быть заполнено")
	}
	if access == "Черновик" {
		return errors.New("нельзя вернуть статус в черновик")
	}

	currentStatus, err := uc.Repository.GetContributorStatus(uint(cid), uint(mid))
	if err != nil {
		return errors.New("ошибка при получении данных о статусе")
	}
	if currentStatus == "Требует подтверждения" {
		err = uc.Repository.UpdateContributorAccessByAdmin(uint(mid), uint(cid), access)
		if err != nil {
			return errors.New("ошибка при обновлении статуса")
		}
		return nil
	}
	return errors.New("невозможно обновить статус заявки")
}
