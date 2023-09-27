package model

import "time"

type Markdown struct {
	Markdown_ID  int `gorm:"primarykey;autoIncrement"`
	Name         string
	Content      string
	Status       string
	Created_Time time.Time `json:"start_date"`
	Moderator_ID uint
}
