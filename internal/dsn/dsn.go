package dsn

import (
	"fmt"
	"log"
)

// Генерируем строку подключения к базе данных
func FromEnv() string {
	// host := os.Getenv("DB_HOST")
	// if host == "" {
	// 	return ""
	// }

	host := "127.0.0.1"

	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// pass := os.Getenv("DB_PASS")
	// dbname := os.Getenv("DB_NAME")

	port := "5450"
	user := "root"
	pass := "S1cParv1sMagna89141403116"
	dbname := "notek_app"

	log.Println(user, port, pass, dbname)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
