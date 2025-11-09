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
	srvproduct "github.com/ftfmtavares/shipping-optimizer/internal/services/product"
	srvshipping "github.com/ftfmtavares/shipping-optimizer/internal/services/shipping"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()
	server := server.NewHTTPServer(server.HTTPServerConfig{
		Address:      cfg.ServerAddress,
		Port:         cfg.ServerPort,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Logger:       instrumentation.NewLogger(),
		Env:          cfg.Env,
	})

	servicesRegistration(ctx, &server)
	server.WithHealthCheck()

	server.StartHTTPServerAsync()
	server.WithShutdownGracefully()
}

func servicesRegistration(ctx context.Context, server *server.HTTPServer) {
	rep := repositories.NewAPIRepositories()

	shippingOptimizer := srvshipping.NewOptimizer(rep.PackSizes)
	server.WithRoute("/product/{pid}/shipping-calculation", api.ShippingCalculation(ctx, shippingOptimizer), http.MethodOptions, http.MethodGet)

	productConfigurator := srvproduct.NewConfigurator(rep.PackSizes)
	server.WithRoute("/product/{pid}/packsizes", api.ProductPackSizes(ctx, productConfigurator), http.MethodOptions, http.MethodGet)
	server.WithRoute("/product/{pid}/packsizes", api.StoreProductPackSizes(ctx, productConfigurator), http.MethodOptions, http.MethodPost)
}
