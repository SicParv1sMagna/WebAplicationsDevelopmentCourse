package usecase

import "project/internal/http/repository"

type UseCase struct {
	Repository *repository.Repository
}

func NewUseCase(r *repository.Repository) *UseCase {
	return &UseCase{
		Repository: r,
	}
}
