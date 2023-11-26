package repository

import (
	"os"
	"project/internal/app/minio"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Структура подключения к БД
type Repository struct {
	db    *gorm.DB
	minio *minio.MinioClient
	redis *redis.Client
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

	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Repository{
		db:    db,
		minio: minio,
		redis: redis,
	}, nil
}
