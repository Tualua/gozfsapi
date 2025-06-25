package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Tualua/gozfsapi/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimw "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	swagger, err := server.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %s", err)
		os.Exit(1)
	}
	defer logger.Sync()
	zfsApi := server.NewApiServer(logger)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	e.Use(oapimw.OapiRequestValidator(swagger))
	server.RegisterHandlers(e, zfsApi)

	go func() {
		logger.Info("Starting server on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", zap.Error(err))
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped gracefully")
}
