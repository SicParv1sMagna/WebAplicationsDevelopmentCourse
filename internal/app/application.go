package app

import (
	"project/internal/dsn"
	"project/internal/repository"

	"github.com/joho/godotenv"
)

type Application struct {
	repository *repository.Repository
}

func New() (Application, error) {
	_ = godotenv.Load()
	repo, err := repository.New(dsn.FromEnv())

	if err != nil {
		return Application{}, err
	}

	return Application{repository: repo}, nil
}
