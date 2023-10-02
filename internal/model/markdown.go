package model

import "time"

// ORM-модель таблицы Markdown из базы данных
type Markdown struct {
	Markdown_ID  int       `json:"Markdown_ID" gorm:"primarykey;autoIncrement"`
	Name         string    `json:"Name"`
	Content      string    `json:"Content"`
	Status       string    `json:"Status"`
	Created_Time time.Time `json:"start_date"`
	Moderator_ID uint
}
