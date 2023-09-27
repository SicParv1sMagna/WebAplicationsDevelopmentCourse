package model

type User struct {
	User_ID    int `gorm:"primarykey;autoIncrement"`
	FirstName  string
	SecondName string
	MiddleName string
	Email      string
	Password   string
}
