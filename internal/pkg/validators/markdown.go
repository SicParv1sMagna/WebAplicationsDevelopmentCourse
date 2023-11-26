package validators

import (
	"errors"
	"project/internal/model"
)

func ValidateMarkdown(md model.Markdown) error {
	if md.Name == "" {
		return errors.New("название должно быть заполнено")
	}

	if len(md.Name) >= 26 {
		return errors.New("название не должно первышать 25 символов")
	}

	return nil
}
