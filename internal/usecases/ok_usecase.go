package usecases

import "context"

type OKUseCaseInterface interface {
	HandleOK(ctx context.Context) string
}

type OKUseCase struct{}

func NewOKUseCase() OKUseCaseInterface {
	return &OKUseCase{}
}

func (u *OKUseCase) HandleOK(ctx context.Context) string {
	return "OK"
}
