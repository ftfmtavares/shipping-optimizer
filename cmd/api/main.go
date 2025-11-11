// Package main for shipping optimizer API
package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ftfmtavares/shipping-optimizer/internal/api"
	"github.com/ftfmtavares/shipping-optimizer/internal/config"
	"github.com/ftfmtavares/shipping-optimizer/internal/instrumentation"
	"github.com/ftfmtavares/shipping-optimizer/internal/repositories"
	"github.com/ftfmtavares/shipping-optimizer/internal/server"
	"github.com/ftfmtavares/shipping-optimizer/internal/services/order"
	"github.com/ftfmtavares/shipping-optimizer/internal/services/product"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()
	server := server.NewHTTPServer(server.HTTPServerConfig{
		Address:      cfg.ServerAddress,
		Port:         cfg.ServerPort,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Logger:       instrumentation.NewLogger(),
	})

	servicesRegistration(ctx, &server)
	staticWeb(&server)
	server.WithHealthCheck()

	server.StartHTTPServerAsync()
	server.WithShutdownGracefully()
}

func servicesRegistration(ctx context.Context, server *server.HTTPServer) {
	rep := repositories.NewAPIRepositories()

	shippingOptimizer := order.NewOptimizer(rep.PackSizes)
	server.WithServiceHandler("/product/{pid}/shipping-calculation", api.OrderCalculation(ctx, shippingOptimizer), http.MethodOptions, http.MethodGet)

	productConfigurator := product.NewConfigurator(rep.PackSizes)
	server.WithServiceHandler("/product/{pid}/packsizes", api.ProductPackSizes(ctx, productConfigurator), http.MethodOptions, http.MethodGet)
	server.WithServiceHandler("/product/{pid}/packsizes", api.StoreProductPackSizes(ctx, productConfigurator), http.MethodOptions, http.MethodPost)
}

func staticWeb(server *server.HTTPServer) {
	server.WithStatic("/", "web")
}
