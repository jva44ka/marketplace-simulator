package app

import (
	"fmt"
	"github.com/jva44ka/ozon-simulator-go/internal/app/handlers/create_review_handler"
	"github.com/jva44ka/ozon-simulator-go/internal/app/handlers/get_reviews_by_sku_handler"
	"github.com/jva44ka/ozon-simulator-go/internal/domain/reviews/repository"
	"github.com/jva44ka/ozon-simulator-go/internal/domain/reviews/service"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/config"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/http/middlewares"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/http/round_trippers"
	"net"
	"net/http"
)

type App struct {
	config *config.Config
	server http.Server
}

func NewApp(configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{
		config: configImpl,
	}

	app.server.Handler = boostrapHandler(configImpl)

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return app.server.Serve(l)
}

func boostrapHandler(config *config.Config) http.Handler {
	tr := http.DefaultTransport
	tr = round_trippers.NewTimerRoundTipper(tr)

	client := http.Client{Transport: tr}

	productService := product_service.NewProductService(
		client,
		config.Products.Token,
		fmt.Sprintf("%s://%s:%s", config.Products.Schema, config.Products.Host, config.Products.Port),
	)

	reviewRepository := repository.NewReviewRepository(100)
	reviewService := service.NewReviewService(reviewRepository, productService)

	mx := http.NewServeMux()
	mx.Handle("POST /products/{sku}/reviews", create_review_handler.NewCreateReviewHandler(reviewService))
	mx.Handle("GET /products/{sku}/reviews", get_reviews_by_sku_handler.NewGetReviewsBySkuHandler(reviewService))

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware
}
