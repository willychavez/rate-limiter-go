package httpserver

import (
	"net/http"

	"github.com/willychavez/rate-limiter-go/internal/infra/web/controllers"
	"github.com/willychavez/rate-limiter-go/internal/infra/web/middleware"
	"github.com/willychavez/rate-limiter-go/internal/usecases"
)

type HTTPServer struct {
	okController          *controllers.OKController
	rateLimiterMiddleware *middleware.RateLimiterMiddleware
}

func NewHTTPServer(rateLimiterUseCase usecases.RateLimiterUseCaseInterface) *HTTPServer {
	// OK UseCase e Controller
	okUseCase := usecases.NewOKUseCase()
	okController := controllers.NewOKController(okUseCase)

	// Middleware
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiterUseCase)

	return &HTTPServer{
		okController:          okController,
		rateLimiterMiddleware: rateLimiterMiddleware,
	}
}

func (s *HTTPServer) Start(addr string) error {
	http.Handle("/", s.rateLimiterMiddleware.Handle(http.HandlerFunc(s.okController.HandleRequest)))
	return http.ListenAndServe(addr, nil)
}
