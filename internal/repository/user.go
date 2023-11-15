package repository

import "project/internal/model"

// Создание пользователя
func (r *Repository) CreateUser(user model.User) error {
	err := r.db.Table("User").Create(&user).Error
	return err
}

// Получаем пользователя по его почте
func (r *Repository) GetUserByEmail(email string) (model.User, error) {
	var user model.User

	// Делаем запрос на получение пользователя из БД
	err := r.db.Table(`"User"`).Where("Email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// Удаление пользователя по ID
func (r *Repository) DeleteUserByID(uID uint) error {
	err := r.db.Table(`"User"`).Where("User_ID = ?", uID).Update("Status", "Удален").Error
	return err
}

func (r *Repository) EditUserData(user model.User, uID uint) error {
	err := r.db.Table(`"User"`).Where("User_ID = ?", uID).Updates(&user).Error
	return err
}

func (r *Repository) GetUserById(id uint) (model.User, error) {
	var user model.User
	err := r.db.Table(`"User"`).Where("User_ID = ?", id).First(&user).Error

	return user, err
}
