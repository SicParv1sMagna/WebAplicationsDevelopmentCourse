package app

import (
	"project/internal/dsn"
	"project/internal/http/delivery"
	"project/internal/http/repository"
	"project/internal/http/usecase"

	"github.com/joho/godotenv"
)

type Application struct {
	repository *repository.Repository
	usecase    *usecase.UseCase
	delivery   *delivery.Delivery
}

func New() (Application, error) {
	_ = godotenv.Load()
	repo, err := repository.New(dsn.FromEnv())
	uc := usecase.NewUseCase(repo)
	d := delivery.NewDelivery(uc)

	if err != nil {
		return Application{}, err
	}

	return Application{
		repository: repo,
		usecase:    uc,
		delivery:   d,
	}, nil
}
