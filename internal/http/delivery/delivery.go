package delivery

import "project/internal/http/usecase"

type Delivery struct {
	usecase *usecase.UseCase
}

func NewDelivery(uc *usecase.UseCase) *Delivery {
	return &Delivery{
		usecase: uc,
	}
}
