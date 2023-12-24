package model

import "time"

// ORM-модель таблицы Markdown из базы данных
type Markdown struct {
	Markdown_ID   int       `json:"Markdown_ID" gorm:"primarykey;autoIncrement"`
	Name          string    `json:"Name" gorm:"column:Name"`
	Content       string    `json:"Content" gorm:"column:Content"`
	Status        string    `json:"Status" gorm:"column:Status"`
	Created_Time  time.Time `json:"start_date"`
	User_ID       uint
	PhotoURL      string `json:"PhotoURL" gorm:"column:PhotoURL"`
	ContributorID int    // Внешний ключ для связи с Contributor
}

type MarkdownWithDates struct {
	Markdown
	Created_Date    *time.Time `json:"created_date"`
	Formed_Date     *time.Time `json:"formed_date"`
	Completion_Date *time.Time `json:"completion_date"`
}
