package model

import "time"

// ORM-модель таблицы Contributor из базы данных
type Contributor struct {
	Contributor_ID  int `gorm:"primarykey;autoIncrement"`
	User_ID         int
	Created_Date    time.Time `json:"start_date"`
	Formed_Date     time.Time `json:"start_date"`
	Completion_Date time.Time `json:"start_date"`
	Status          string
}
