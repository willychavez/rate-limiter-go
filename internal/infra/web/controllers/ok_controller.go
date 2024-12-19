package controllers

import (
	"net/http"

	"github.com/willychavez/rate-limiter-go/internal/usecases"
)

type OKController struct {
	okUseCase usecases.OKUseCaseInterface
}

func NewOKController(okUseCase usecases.OKUseCaseInterface) *OKController {
	return &OKController{okUseCase: okUseCase}
}

func (c *OKController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := c.okUseCase.HandleOK(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
