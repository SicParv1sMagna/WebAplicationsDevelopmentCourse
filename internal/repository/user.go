package repository

import "project/internal/model"

// Создание пользователя
func (r *Repository) CreateUser(user model.User) error {
	// Строка добавления пользователя с заданными данными в БД
	sql := `INSERT INTO "User" (First_Name, Second_Name, Middle_Name, Email, "Password") VALUES (?, ?, ?, ?, ?)`

	// Добавление пользователя в базу данных
	err := r.db.Exec(sql, user.FirstName, user.SecondName, user.MiddleName, user.Email, user.Password).Error
	if err != nil {
		return err
	}

	return nil
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
func (r *Repository) DeleteUserByID() {
}
