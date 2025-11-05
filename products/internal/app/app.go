package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/jva44ka/ozon-simulator-go/internal/app/handlers/get_product_by_sku_handler"
	"github.com/jva44ka/ozon-simulator-go/internal/domain/repository"
	"github.com/jva44ka/ozon-simulator-go/internal/domain/service"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/config"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/http/middlewares"
	"github.com/jva44ka/ozon-simulator-go/internal/infra/http/round_trippers"
	httpSwagger "github.com/swaggo/http-swagger"
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

	productRepository := repository.NewProductRepository(100)
	productService := service.NewProductService(productRepository)

	mx := http.NewServeMux()
	mx.Handle("GET /products/{sku}", get_product_by_sku_handler.NewGetProductsBySkuHandler(productService))
	mx.Handle("GET /swagger/*", httpSwagger.WrapHandler)

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware
}
