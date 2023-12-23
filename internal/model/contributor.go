package model

import (
	"time"
)

// ORM-модель таблицы Contributor из базы данных
type Contributor struct {
	Contributor_ID  int `gorm:"primarykey;autoIncrement;column:contributor_id"`
	User_ID         int
	Created_Date    *time.Time `json:"created_date"`
	Formed_Date     *time.Time `json:"formed_date"`
	Completion_Date *time.Time `json:"completion_date"`
	Email           string     `json:"email" gorm:"column:email"`
}

type ContributorWithMarkdowns struct {
	Contributor
	Markdown []Markdown `gorm:"foreignKey:contributor_id;references:contributor_id"`
}
