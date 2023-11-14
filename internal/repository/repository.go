package repository

import (
	"project/internal/app/minio"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Структура подключения к БД
type Repository struct {
	db    *gorm.DB
	minio *minio.MinioClient
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minio, err := minio.NewMinioClient()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:    db,
		minio: minio,
	}, nil
}
