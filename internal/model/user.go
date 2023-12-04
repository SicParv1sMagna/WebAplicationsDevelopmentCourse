package model

import "project/internal/pkg/roles"

// ORM-модель таблицы User из базы данных
type User struct {
	User_ID        uint64     `json:"User_ID" gorm:"primarykey;autoIncrement"`
	FirstName      string     `json:"FirstName"`
	SecondName     string     `json:"SecondName"`
	MiddleName     string     `json:"MiddleName"`
	Email          string     `json:"Email"`
	Password       string     `json:"Password" gorm:"column:Password"`
	RepeatPassword string     `json:"RepeatPassword" gorm:"-"`
	Role           roles.Role `json:"Role" gorm:"column:Role"`
}

type UserRegisterReq struct {
	FirstName      string `json:"FirstName"`
	SecondName     string `json:"SecondName"`
	MiddleName     string `json:"MiddleName"`
	Email          string `json:"Email"`
	Password       string `json:"Password" gorm:"column:Password"`
	RepeatPassword string `json:"RepeatPassword" gorm:"-"`
}

type UserAuthReq struct {
	Email    string `json:"Email"`
	Password string `json:"Password" gorm:"column:Password"`
}
