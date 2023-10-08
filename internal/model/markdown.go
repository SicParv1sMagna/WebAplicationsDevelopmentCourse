package model

import "time"

// ORM-модель таблицы Markdown из базы данных
type Markdown struct {
	Markdown_ID  int       `json:"Markdown_ID" gorm:"primarykey;autoIncrement"`
	Name         string    `json:"Name" gorm:"column:Name"`
	Content      string    `json:"Content" gorm:"column:Content"`
	Status       string    `json:"Status" gorm:"column:Status"`
	Created_Time time.Time `json:"start_date"`
	User_ID      uint
}
