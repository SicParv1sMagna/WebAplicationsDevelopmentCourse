package model

// ORM-модель таблицы User из базы данных
type User struct {
	User_ID        uint64 `json:"User_ID" gorm:"primarykey;autoIncrement"`
	FirstName      string `json:"FirstName"`
	SecondName     string `json:"SecondName"`
	MiddleName     string `json:"MiddleName"`
	Email          string `json:"Email"`
	Password       string `json:"Password" gorm:"column:Password"`
	RepeatPassword string `json:"RepeatPassword" gorm:"-"`
}
